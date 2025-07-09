package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/server"
)

var (
	_ resource.Resource                = &ServerResource{}
	_ resource.ResourceWithConfigure   = &ServerResource{}
	_ resource.ResourceWithImportState = &ServerResource{}
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	client *service.ApiServerService
}

func (r *ServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.ServerResourceSchema(ctx)
}

func (r *ServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetServerService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	serverService, ok := client.(*service.ApiServerService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = serverService
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Missing required field",
			"'name' field is required when creating a server",
		)
	}

	if data.ApplianceID.IsNull() || data.ApplianceID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("appliance_id"),
			"Missing required field",
			"'appliance_id' field is required when creating a server",
		)
	}

	if (data.DatacenterID.IsNull() || data.DatacenterID.ValueString() == "") &&
		(data.SiteID.IsNull() || data.SiteID.ValueString() == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("datacenter_id"),
			"Missing required field",
			"Either 'datacenter_id' or 'site_id' field is required when creating a server",
		)
	}

	if data.Hardware.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hardware"),
			"Missing required field",
			"'hardware' field is required when creating a server",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := data.ToCreateRequest()

	tflog.Info(ctx, fmt.Sprintf("Creating server: %s", createRequest.Name))

	apiResponse, err := r.client.CreateServer(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			fmt.Sprintf("Could not create server: %s", err),
		)
		return
	}

	model, diags := models.NewServerResourceModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created server with ID: %s", model.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading server with ID: %s", id))

	apiResponse, err := r.client.GetServer(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Server with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading server",
			fmt.Sprintf("Could not read server: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Server with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	readModel, diags := models.NewServerResourceModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ServerResourceModel
	var state models.ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating server with ID: %s", id))

	updateRequest := plan.ToUpdateRequest()

	updatedServer, err := r.client.UpdateServer(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating server",
			fmt.Sprintf("Could not update server: %s", err),
		)
		return
	}

	updatedModel, diags := models.NewServerResourceModel(ctx, updatedServer)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully updated server with ID: %s", id))

	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting server with ID: %s", id))

	err := r.client.DeleteServer(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Server %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting server",
			fmt.Sprintf("Could not delete server: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted server with ID: %s", id))
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
