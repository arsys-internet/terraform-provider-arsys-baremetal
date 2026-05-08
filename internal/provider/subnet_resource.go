package provider

import (
	"context"
	"errors"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/subnet"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &SubnetResource{}
	_ resource.ResourceWithConfigure = &SubnetResource{}
)

func NewSubnetResource() resource.Resource {
	return &SubnetResource{}
}

type SubnetResource struct {
	client *service.ApiSubnetService
}

func (r *SubnetResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (r *SubnetResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.SubnetResourceSchema(ctx)
}

func (r *SubnetResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetSubnetService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	subnetService, ok := client.(*service.ApiSubnetService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = subnetService
}

func (r *SubnetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.SubnetModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := data.ToCreateRequest()

	apiResponse, err := r.client.CreateSubnet(&createRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating subnet",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while creating subnet. Please report this issue to the provider developers.",
		)
		return
	}

	finalModel, diags := models.NewSubnetModelFromResponse(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to create final resource model")
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *SubnetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading subnet with ID: %s", id))

	_, err := r.client.GetSubnet(id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			tflog.Info(ctx, fmt.Sprintf("Subnet %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading subnet",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *SubnetResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource does not support updates. Check the provider documentation for more details.",
	)
}

func (r *SubnetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.SubnetModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	err := r.client.DeleteSubnet(id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			tflog.Info(ctx, fmt.Sprintf("Subnet %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting subnet",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}
}
