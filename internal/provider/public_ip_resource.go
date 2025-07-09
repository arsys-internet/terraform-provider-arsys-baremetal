package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicIp"
	"terraform-provider-arsys-baremetal/internal/util"
	"time"
)

var (
	_ resource.Resource                = &PublicIpResource{}
	_ resource.ResourceWithConfigure   = &PublicIpResource{}
	_ resource.ResourceWithImportState = &PublicIpResource{}
)

func NewPublicIpResource() resource.Resource {
	return &PublicIpResource{}
}

type PublicIpResource struct {
	client *service.ApiPublicIpService
}

func (r *PublicIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (r *PublicIpResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PublicIpResourceSchema(ctx)
}

func (r *PublicIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicIpService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicIpService, ok := client.(*service.ApiPublicIpService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = publicIpService
}

func (r *PublicIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PublicIpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading public ip with ID: %s", id))

	apiResponse, err := r.client.GetPublicIp(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Public ip with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading public ip",
			fmt.Sprintf("Could not read public ip: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Public ip with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	readModel, diags := models.NewPublicIpResourceModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)

}

func (r *PublicIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PublicIpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := data.ToCreateRequest()

	apiResponse, err := r.client.CreatePublicIp(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating public ip",
			fmt.Sprintf("Could not create public ip: %s", err),
		)
		return
	}

	model, diags := models.NewPublicIpResourceModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created public ip with ID: %s", model.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *PublicIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PublicIpResourceModel
	var state models.PublicIpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating public ip with ID: %s", id))

	updateRequest := plan.ToUpdateRequest()

	updatedPublicIp, err := r.client.UpdatePublicIp(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating public ip",
			fmt.Sprintf("Could not update public ip: %s", err),
		)
		return
	}

	updatedModel, diags := models.NewPublicIpResourceModel(ctx, updatedPublicIp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully updated public ip with ID: %s", id))

	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
}

func (r *PublicIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PublicIpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting public ip with ID: %s", id))

	err := r.client.DeletePublicIp(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Public ip %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting public ip",
			fmt.Sprintf("Could not delete public ip: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getIpTimeout()

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

	tflog.Info(ctx, fmt.Sprintf("Deleted public ip with ID: %s", id))
}

func (r *PublicIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getIpTimeout() (time.Duration, time.Duration, time.Duration) {
	var timeout = util.GetTimeoutFromEnv("IP_DEFAULT_TIMEOUT", time.Minute)
	var retryInterval = util.GetTimeoutFromEnv("IP_DEFAULT_RETRY_INTERVAL", time.Second)
	var minTimeout = util.GetTimeoutFromEnv("IP_DEFAULT_MIN_TIMEOUT", time.Second)

	return timeout, retryInterval, minTimeout
}
