package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicnetwork"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &PublicNetworkIpsResource{}
	_ resource.ResourceWithConfigure   = &PublicNetworkIpsResource{}
	_ resource.ResourceWithImportState = &PublicNetworkIpsResource{}
)

func NewPublicNetworkIpsResource() resource.Resource {
	return &PublicNetworkIpsResource{}
}

type PublicNetworkIpsResource struct {
	client *service.ApiPublicNetworkIpService
}

func (r *PublicNetworkIpsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network_ips"
}

func (r *PublicNetworkIpsResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PublicNetworkIpResourceSchema(ctx)
}

func (r *PublicNetworkIpsResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicNetworkIpService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	PublicNetworkIpService, ok := client.(*service.ApiPublicNetworkIpService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = PublicNetworkIpService
}

func (r *PublicNetworkIpsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PublicNetworkIpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetworkId := data.PublicNetworkId.ValueString()
	if publicNetworkId == "" {
		publicNetworkId = data.Id.ValueString()
	}

	if publicNetworkId == "" {
		resp.Diagnostics.AddError(
			"Missing Public Network ID",
			"Could not determine the public network ID from state. Please re-import the resource with a valid ID.",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading IPs for public network with ID: %s", publicNetworkId))

	apiResponse, err := r.client.GetPublicNetworkIps(publicNetworkId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Public network %s not found, removing from state", publicNetworkId))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading public network IPs",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving public network IPs. Please try again or report this issue to the provider developers",
		)
		return
	}

	data.PublicNetworkId = types.StringValue(publicNetworkId)
	data.Id = types.StringValue(publicNetworkId)

	if len(data.Ips) == 0 {
		ipIds := make([]types.String, len(apiResponse))
		for i, ip := range apiResponse {
			ipIds[i] = types.StringValue(ip.Id)
		}
		data.Ips = ipIds
		data.Action = types.BoolValue(true)
	}

	readModel, diags := models.NewPublicNetworkIpResourceModel(ctx, &data, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *PublicNetworkIpsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PublicNetworkIpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetworkId := data.PublicNetworkId.ValueString()
	isAssigning := data.Action.ValueBool()

	createRequest := data.ToCreateRequest()

	apiResponse, err := r.client.AssignIpToPublicNetwork(publicNetworkId, &createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error processing IP assignment/unassignment",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while assigning IPs to public network. Please try again or report this issue to the provider developers",
		)
		return
	}

	if data.Id.IsNull() || data.Id.ValueString() == "" {
		data.Id = types.StringValue(publicNetworkId)
	}

	if isAssigning {

		timeouts := util.GetResourceTimeouts("PUBLIC_NETWORK")
		waitOptions := util.NewWaitOptions(
			timeouts.Default,
			timeouts.RetryInterval,
			timeouts.MinTimeout,
			[]string{util.StatusConfiguring},
			[]string{util.StatePoweredOn, util.StatePoweredOff},
		)

		for _, ip := range *apiResponse {
			compositeId := service.CreateCompositeID(publicNetworkId, ip.Id)

			_, diags := util.WaitForResourceState(
				ctx,
				compositeId,
				r.client,
				waitOptions,
			)

			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				resp.Diagnostics.AddError(
					"Wait failed for IP assignment",
					fmt.Sprintf("IP %s failed to reach ready state", ip.Id),
				)
				return
			}
		}

		tflog.Info(ctx, fmt.Sprintf("All %d IPs are now ready", len(*apiResponse)))
	} else {
		tflog.Info(ctx, fmt.Sprintf("Unassigned %d IPs from public network", len(*apiResponse)))
	}

	finalIps, err := r.client.GetPublicNetworkIps(publicNetworkId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting final IP state",
			fmt.Sprintf("Could not get final IP state: %s", err.Error()),
		)
		return
	}

	finalModel, diags := models.NewPublicNetworkIpResourceModel(ctx, &data, finalIps)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if isAssigning {
		tflog.Info(ctx, fmt.Sprintf("Successfully assigned IPs to public network %s", publicNetworkId))
	} else {
		tflog.Info(ctx, fmt.Sprintf("Successfully unassigned IPs from public network %s", publicNetworkId))
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)

}

func (r *PublicNetworkIpsResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource cannot be updated. Please check your Terraform configuration.",
	)
}

func (r *PublicNetworkIpsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PublicNetworkIpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetworkId := data.PublicNetworkId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Removing public network IP assignment from Terraform state only: %s", publicNetworkId))
	tflog.Info(ctx, "Note: Server IPs remain assigned to the public network - only removing from Terraform management")
}

func (r *PublicNetworkIpsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
