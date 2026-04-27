package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/sshkey"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &SshKeyResource{}
	_ resource.ResourceWithConfigure   = &SshKeyResource{}
	_ resource.ResourceWithImportState = &SshKeyResource{}
)

func NewSshKeyResource() resource.Resource {
	return &SshKeyResource{}
}

type SshKeyResource struct {
	client *service.ApiSshKeyService
}

func (r *SshKeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (r *SshKeyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.SshKeyResourceSchema(ctx)
}

func (r *SshKeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetSshKeyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	sshKeyService, ok := client.(*service.ApiSshKeyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = sshKeyService
}

func (r *SshKeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.SshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading SSH key with ID: %s", id))

	apiResponse, err := r.client.GetSshKey(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("SSH key with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading SSH key",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving SSH key. Please try again or report this issue to the provider developers",
		)
		return
	}

	readModel, diags := models.NewSshKeyModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)

}

func (r *SshKeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.SshKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.IsUnknown() || data.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Missing required field",
			"Either 'name' field is required when creating a SSH key",
		)
	}

	createRequest := data.ToCreateRequest()

	apiResponse, err := r.client.CreateSshKey(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SSH key",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while creating SSH key. Please try again or report this issue to the provider developers",
		)
		return
	}

	model, diags := models.NewSshKeyModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created SSH key with Id: %s", model.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *SshKeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.SshKeyModel
	var state models.SshKeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating SSH key with ID: %s", id))

	if plan.Name.IsNull() || plan.Name.IsUnknown() || plan.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Missing required field",
			"Either 'name' field is required when updating a SSH key",
		)
	}

	updateRequest := plan.ToUpdateRequest()

	updatedSshKey, err := r.client.UpdateSshKey(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SSH key",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if updatedSshKey == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while updating SSH key. Please try again or report this issue to the provider developers",
		)
		return
	}

	updatedModel, diags := models.NewSshKeyModel(ctx, updatedSshKey)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully updated SSH key with ID: %s", id))

	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
}

func (r *SshKeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.SshKeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting SSH key with ID: %s", id))

	err := r.client.DeleteSshKey(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("SSH key %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting SSH key",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted SSH key with ID: %s", id))
}

func (r *SshKeyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
