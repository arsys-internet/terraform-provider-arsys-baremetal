package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/privateNetwork"
)

var (
	_ resource.Resource                = &PrivateNetworkResource{}
	_ resource.ResourceWithConfigure   = &PrivateNetworkResource{}
	_ resource.ResourceWithImportState = &PrivateNetworkResource{}
)

func NewPrivateNetworkResource() resource.Resource {
	return &PrivateNetworkResource{}
}

type PrivateNetworkResource struct {
	client *service.ApiPrivateNetworkService
}

func (r *PrivateNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_network"
}

func (r *PrivateNetworkResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PrivateNetworkResourceSchema(ctx)
}

func (r *PrivateNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPrivateNetworkService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	privateNetworkService, ok := client.(*service.ApiPrivateNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = privateNetworkService
}

func (r *PrivateNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PrivateNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Missing required field",
			"'name' field is required when creating a private network",
		)
	}

	if data.DatacenterID.IsNull() || data.DatacenterID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("datacenter_id"),
			"Missing required field",
			"'datacenter_id' field is required when creating a private network",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := data.ToCreateRequest()

	tflog.Info(ctx, fmt.Sprintf("Creating private network: %s", createRequest.Name))

	apiResponse, err := r.client.CreatePrivateNetwork(&createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating private network",
			fmt.Sprintf("Could not create private network: %s", err),
		)
		return
	}

	model, diags := models.NewPrivateNetworkResource(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created private network with ID: %s", model.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *PrivateNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PrivateNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading private network with ID: %s", id))

	apiResponse, err := r.client.GetPrivateNetwork(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Private network Not Found",
				fmt.Sprintf("Private network with id %s not found", id),
			)
			tflog.Info(ctx, fmt.Sprintf("Private network with ID %s not found", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error reading private network",
			fmt.Sprintf("Could not read private network: %s", err),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("Private network not found"),
		)
		return
	}

	readModel, diags := models.NewPrivateNetworkResource(ctx, apiResponse)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *PrivateNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.PrivateNetworkResourceModel
	var state models.PrivateNetworkResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating private network with ID: %s", id))

	updateRequest := plan.ToUpdateRequest()

	updatedPrivateNetwork, err := r.client.UpdatePrivateNetwork(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating private network",
			fmt.Sprintf("Could not update private network: %s", err),
		)
		return
	}

	updatedModel, diags := models.NewPrivateNetworkResource(ctx, updatedPrivateNetwork)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully updated private network with ID: %s", id))

	diags = resp.State.Set(ctx, updatedModel)
	resp.Diagnostics.Append(diags...)
}

func (r *PrivateNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PrivateNetworkResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting private network with ID: %s", id))

	err := r.client.DeletePrivateNetwork(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Private network %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting private network",
			fmt.Sprintf("Could not delete private network: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted private network with ID: %s", id))
}

func (r *PrivateNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
