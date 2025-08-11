package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicNetwork"
	"terraform-provider-arsys-baremetal/internal/util"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &PublicNetworkResource{}
	_ resource.ResourceWithConfigure   = &PublicNetworkResource{}
	_ resource.ResourceWithImportState = &PublicNetworkResource{}
)

func NewPublicNetworkResource() resource.Resource {
	return &PublicNetworkResource{}
}

type PublicNetworkResource struct {
	client *service.ApiPublicNetworkService
}

func (r *PublicNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network"
}

func (r *PublicNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PublicNetworkResourceSchema(ctx)
}

func (r *PublicNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicNetworkService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicNetworkService, ok := client.(*service.ApiPublicNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = publicNetworkService
}

func (r *PublicNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PublicNetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading public network with ID: %s", id))

	apiResponse, err := r.client.GetPublicNetwork(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Public network with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading public network",
			fmt.Sprintf("Could not read public network: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Public network with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	readModel, diags := models.NewPublicNetworkModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)

}

func (r *PublicNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PublicNetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.DatacenterId.IsNull() || data.DatacenterId.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("datacenter_id"),
			"Missing required field",
			"Either 'datacenter_id' field is required when creating a public network",
		)
	}

	createRequest := data.ToCreateRequest()

	apiResponse, err := r.client.CreatePublicNetwork(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating public network",
			fmt.Sprintf("Could not create public network: %s", err),
		)
		return
	}

	model, diags := models.NewPublicNetworkModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getPublicNetworkTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
		[]string{util.StateDeploying},
		[]string{util.StatePoweredOn, util.StatePoweredOff},
	)

	waitResult, diags := util.WaitForResourceState(
		ctx,
		apiResponse.ID,
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	finalModel := model
	if waitResult != nil && waitResult.Resource != nil {
		if publicNetworkModel, ok := waitResult.Resource.(*models.PublicNetworkModel); ok {
			finalModel = publicNetworkModel
			tflog.Info(ctx, fmt.Sprintf("Public network reached final state: %s", waitResult.FinalState))
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Created public network with ID: %s", finalModel.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *PublicNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PublicNetworkModel
	var state models.PublicNetworkModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating public network with ID: %s", id))

	updateRequest := plan.ToUpdateRequest()

	updatedPublicNetwork, err := r.client.UpdatePublicNetwork(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating public network",
			fmt.Sprintf("Could not update public network: %s", err),
		)
		return
	}

	updatedModel, diags := models.NewPublicNetworkModel(ctx, updatedPublicNetwork)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully updated public network with ID: %s", id))

	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
}

func (r *PublicNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PublicNetworkModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting public network with ID: %s", id))

	err := r.client.DeletePublicNetwork(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Public network %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting public network",
			fmt.Sprintf("Could not delete public network: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getPublicNetworkTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
		[]string{util.StateRemoving},
		[]string{util.StateDeleted},
	)

	_, diags := util.WaitForResourceState(
		ctx,
		data.ID.ValueString(),
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted public network with ID: %s", id))
}

func (r *PublicNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getPublicNetworkTimeout() (time.Duration, time.Duration, time.Duration) {
	var timeout = util.GetTimeoutFromEnv("PUBLIC_NETWORK_DEFAULT_TIMEOUT", time.Minute)
	var retryInterval = util.GetTimeoutFromEnv("PUBLIC_NETWORK_DEFAULT_RETRY_INTERVAL", time.Second)
	var minTimeout = util.GetTimeoutFromEnv("PUBLIC_NETWORK_DEFAULT_MIN_TIMEOUT", time.Second)

	return timeout, retryInterval, minTimeout
}
