package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/private_network"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PrivateNetworkServerRemoveResource{}

func NewPrivateNetworkServerRemoveResource() resource.Resource {
	return &PrivateNetworkServerRemoveResource{}
}

type PrivateNetworkServerRemoveResource struct {
	client *service.ApiPrivateNetworkService
}

func (r *PrivateNetworkServerRemoveResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_network_server_remove"
}

func (r *PrivateNetworkServerRemoveResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PrivateNetworkServerResourceRemoveSchema(ctx)
}

func (r *PrivateNetworkServerRemoveResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	networkService, ok := client.(*service.ApiPrivateNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = networkService
}

func (r *PrivateNetworkServerRemoveResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PrivateNetworkServerRemoveModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateNetworkId := data.Id.ValueString()
	serverId := data.ServerId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Removing server %s from private network %s", serverId, privateNetworkId))

	_, err := r.client.DeletePrivateNetworkServer(privateNetworkId, serverId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Server %s or private network %s not found - may already be deleted", serverId, privateNetworkId))
		} else {
			resp.Diagnostics.AddError(
				"Error removing server from private network",
				fmt.Sprintf("Error: %s", err.Error()),
			)
			return
		}
	}

	timeouts := util.GetResourceTimeouts("PRIVATE_NETWORKS_OPERATIONS")
	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StateActive},
	)

	_, diags := util.WaitForResourceState(ctx, privateNetworkId, r.client, waitOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Wait for private network state failed")
		return
	}

	privateNetwork, err := r.client.GetPrivateNetwork(privateNetworkId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting final private network state",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if privateNetwork == nil {
		resp.Diagnostics.AddError(
			"Unexpected Private Network State",
			"API returned no Private network data after rule removal",
		)
		return
	}

	finalModel, diags := models.NewPrivateNetworkServerRemoveModel(ctx, serverId, *privateNetwork)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully removed server %s from private network %s", serverId, privateNetworkId))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *PrivateNetworkServerRemoveResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PrivateNetworkServerRemoveModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateNetwork, err := r.client.GetPrivateNetwork(data.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading private network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if privateNetwork == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	updatedModel, diags := models.NewPrivateNetworkServerRemoveModel(
		ctx,
		data.ServerId.ValueString(),
		*privateNetwork,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *PrivateNetworkServerRemoveResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource cannot be updated. Please check your Terraform configuration.",
	)
}

func (r *PrivateNetworkServerRemoveResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Info(ctx, "Removing server removal tracking from Terraform state")
}
