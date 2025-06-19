package util

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"time"
)

// Configuración por defecto para rate limiting
const (
	DefaultRateLimitMultiplier = 3
	RateLimitErrorSubstring    = "429"
)

type ResourceModel interface {
	GetState() string
}

type ResourceService interface {
	GetResource(id string) (ResourceModel, error)
}

type WaitOptions struct {
	Message              string
	Timeout              time.Duration
	RetryInterval        time.Duration
	MinTimeout           time.Duration
	PendingStates        []string
	TargetStates         []string
	RateLimitMultiplier  int
	IgnoreNotFoundErrors bool
}

func NewWaitOptions(timeout, retryInterval, minTimeout time.Duration, pendingStates, targetStates []string) WaitOptions {
	return WaitOptions{
		Timeout:              timeout * time.Minute,
		RetryInterval:        retryInterval,
		MinTimeout:           minTimeout,
		PendingStates:        pendingStates,
		TargetStates:         targetStates,
		RateLimitMultiplier:  DefaultRateLimitMultiplier,
		IgnoreNotFoundErrors: false,
	}
}

func NewWaitOptionsWithDefaults(opts WaitOptions) WaitOptions {
	if opts.RateLimitMultiplier <= 0 {
		opts.RateLimitMultiplier = DefaultRateLimitMultiplier
	}
	if opts.RetryInterval == 0 {
		opts.RetryInterval = 5 * time.Second
	}
	if opts.MinTimeout == 0 {
		opts.MinTimeout = 1 * time.Second
	}
	return opts
}

type WaitResult struct {
	Resource   ResourceModel
	FinalState string
	Duration   time.Duration
}

func WaitForResourceState(ctx context.Context, resourceID string, service ResourceService, options WaitOptions) (*WaitResult, diag.Diagnostics) {
	var diags diag.Diagnostics
	startTime := time.Now()

	tflog.Info(ctx, "Starting wait for resource state", map[string]interface{}{
		"resource_id":    resourceID,
		"target_states":  options.TargetStates,
		"pending_states": options.PendingStates,
		"timeout":        options.Timeout.String(),
		"retry_interval": options.RetryInterval.String(),
	})

	isDeleting := isDeletionOperation(options.PendingStates, options.TargetStates)

	// Verificación inicial del recurso
	resource, err := service.GetResource(resourceID)
	if err != nil && (!isDeleting || !options.IgnoreNotFoundErrors) {
		diags.AddError(
			"Error getting resource",
			fmt.Sprintf("Failed to get resource %s: %s", resourceID, err.Error()),
		)
		return nil, diags
	}

	// Caso especial: recurso no encontrado durante eliminación
	if resource == nil && isDeleting {
		tflog.Info(ctx, "Resource not found during deletion - considering as deleted")
		return &WaitResult{
			Resource:   nil,
			FinalState: StateDeleted,
			Duration:   time.Since(startTime),
		}, diags
	}

	if resource != nil {
		currentState := resource.GetState()
		if contains(options.TargetStates, currentState) {
			tflog.Info(ctx, "Resource already in target state", map[string]interface{}{
				"resource_id": resourceID,
				"state":       currentState,
			})
			return &WaitResult{
				Resource:   resource,
				FinalState: currentState,
				Duration:   time.Since(startTime),
			}, diags
		}
	}

	result, pollDiags := pollResourceState(ctx, resourceID, service, options, isDeleting)
	diags.Append(pollDiags...)

	if result != nil {
		result.Duration = time.Since(startTime)
		tflog.Info(ctx, "Wait operation completed", map[string]interface{}{
			"resource_id":    resourceID,
			"final_state":    result.FinalState,
			"total_duration": result.Duration.String(),
		})
	}

	return result, diags
}

func pollResourceState(ctx context.Context, resourceID string, service ResourceService, options WaitOptions, isDeleting bool) (*WaitResult, diag.Diagnostics) {
	var diags diag.Diagnostics

	ticker := time.NewTicker(options.RetryInterval)
	defer ticker.Stop()

	timeoutCh := time.After(options.Timeout)
	retryCount := 0

	for {
		select {
		case <-ctx.Done():
			diags.AddError(
				"Context canceled",
				"Wait operation was canceled",
			)
			return nil, diags

		case <-timeoutCh:
			diags.AddError(
				"Timeout waiting for resource state",
				fmt.Sprintf("Resource %s did not reach target states %v within timeout %v after %d retries",
					resourceID, options.TargetStates, options.Timeout, retryCount),
			)
			return nil, diags

		case <-ticker.C:
			retryCount++
			result, shouldContinue, pollDiags := checkResourceState(ctx, resourceID, service, options, isDeleting, retryCount)
			diags.Append(pollDiags...)

			if diags.HasError() {
				return nil, diags
			}

			if !shouldContinue {
				return result, diags
			}

			if pollDiags.HasError() && isRateLimitError(pollDiags) {
				adjustedInterval := time.Duration(options.RateLimitMultiplier) * options.RetryInterval
				tflog.Warn(ctx, "Rate limited, adjusting retry interval", map[string]interface{}{
					"original_interval": options.RetryInterval.String(),
					"adjusted_interval": adjustedInterval.String(),
				})
				time.Sleep(adjustedInterval)
			}
		}
	}
}

func checkResourceState(ctx context.Context, resourceID string, service ResourceService, options WaitOptions, isDeleting bool, retryCount int) (*WaitResult, bool, diag.Diagnostics) {
	currentResource, err := service.GetResource(resourceID)

	if err != nil && isRateLimitingError(err) {
		tflog.Warn(ctx, "Rate limited, will retry", map[string]interface{}{
			"retry_count": retryCount,
			"error":       err.Error(),
		})
		return nil, true, diag.Diagnostics{}
	}

	if isDeleting {
		return handleDeletionState(ctx, resourceID, currentResource, err, options)
	}

	return handleNormalState(ctx, resourceID, currentResource, err, options)
}

func handleDeletionState(ctx context.Context, resourceID string, resource ResourceModel, err error, options WaitOptions) (*WaitResult, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err != nil || resource == nil {
		tflog.Info(ctx, "Resource deleted successfully", map[string]interface{}{
			"resource_id": resourceID,
		})
		return &WaitResult{
			Resource:   nil,
			FinalState: StateDeleted,
		}, false, diags
	}

	resourceModel, ok := resource.(ResourceModel)
	if !ok {
		diags.AddError(
			"Type assertion error",
			"Resource does not implement ResourceModel interface",
		)
		return nil, false, diags
	}

	state := resourceModel.GetState()
	if state == "" {
		tflog.Info(ctx, "Resource state empty - considering as deleted")
		return &WaitResult{
			Resource:   nil,
			FinalState: StateDeleted,
		}, false, diags
	}

	return evaluateState(ctx, resourceID, resource, state, options)
}

func handleNormalState(ctx context.Context, resourceID string, resource ResourceModel, err error, options WaitOptions) (*WaitResult, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if err != nil {
		diags.AddError(
			"Error checking resource",
			fmt.Sprintf("Failed to check resource %s: %s", resourceID, err.Error()),
		)
		return nil, false, diags
	}

	if resource == nil {
		diags.AddError(
			"Resource not found",
			fmt.Sprintf("Resource %s not found", resourceID),
		)
		return nil, false, diags
	}

	resourceModel, ok := resource.(ResourceModel)
	if !ok {
		diags.AddError(
			"Type assertion error",
			"Resource does not implement ResourceModel interface",
		)
		return nil, false, diags
	}

	currentState := resourceModel.GetState()
	return evaluateState(ctx, resourceID, resource, currentState, options)
}

func evaluateState(ctx context.Context, resourceID string, resource ResourceModel, currentState string, options WaitOptions) (*WaitResult, bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	if contains(options.TargetStates, currentState) {
		tflog.Info(ctx, "Resource reached target state", map[string]interface{}{
			"resource_id": resourceID,
			"state":       currentState,
		})
		return &WaitResult{
			Resource:   resource,
			FinalState: currentState,
		}, false, diags
	}

	if !contains(options.PendingStates, currentState) {
		diags.AddError(
			"Unexpected resource state",
			fmt.Sprintf("Resource %s is in unexpected state: %s. Expected one of: %v or %v",
				resourceID, currentState, options.PendingStates, options.TargetStates),
		)
		return nil, false, diags
	}

	tflog.Debug(ctx, "Resource still in pending state", map[string]interface{}{
		"resource_id": resourceID,
		"state":       currentState,
	})

	return nil, true, diags // Continuar polling
}

// Funciones utilitarias
func isDeletionOperation(pendingStates, targetStates []string) bool {
	return (contains(pendingStates, StateRemoving) || contains(pendingStates, StateRemoving)) &&
		contains(targetStates, StateDeleted)
}

func isRateLimitingError(err error) bool {
	return err != nil && strings.Contains(err.Error(), RateLimitErrorSubstring)
}

func isRateLimitError(diags diag.Diagnostics) bool {
	for _, d := range diags {
		if strings.Contains(d.Detail(), RateLimitErrorSubstring) {
			return true
		}
	}
	return false
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
