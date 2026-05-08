package provider

import (
	"context"
	"errors"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/privatenetwork"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &PrivateNetworkServersResource{}

func NewPrivateNetworkServersAssignResource() resource.Resource {
	return &PrivateNetworkServersResource{}
}

type PrivateNetworkServersResource struct {
	client *service.ApiPrivateNetworkService
}

func (r *PrivateNetworkServersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_network_servers_assign"
}

func (r *PrivateNetworkServersResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PrivateNetworkServerAssignResourceSchema(ctx)
}

func (r *PrivateNetworkServersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PrivateNetworkServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PrivateNetworkServerAssignModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateNetworkId := data.Id.ValueString()

	assignRequest, diags := data.ToAssignRequest(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResponse, assignErr := r.client.CreatePrivateNetworkServers(privateNetworkId, assignRequest)
	if assignErr != nil {
		resp.Diagnostics.AddError(
			"Error assigning servers to private network",
			fmt.Sprintf("Error: %s", assignErr.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("PRIVATE_NETWORKS_OPERATIONS")
	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StateActive},
	)

	id := apiResponse.Id

	_, waitDiags := util.WaitForResourceState(
		ctx,
		id,
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(waitDiags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Wait for private network state failed")
		return
	}

	privateNetwork, fwErr := r.client.GetPrivateNetwork(id)
	if fwErr != nil {
		resp.Diagnostics.AddError(
			"Error getting final private network state",
			fmt.Sprintf("Could not get final private network state: %s", fwErr.Error()),
		)
		return
	}

	if privateNetwork == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving private network after assign server. Please report this issue to the provider developers.",
		)
		return
	}

	finalModel, diags := models.NewPrivateNetworkServerAssignModel(ctx, *privateNetwork)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully assigned %d servers to private network: %s", len(assignRequest.Servers), privateNetworkId))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *PrivateNetworkServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PrivateNetworkServerAssignModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyId := data.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading private network assignment: %s", firewallPolicyId))

	apiResponse, err := r.client.GetPrivateNetwork(firewallPolicyId)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			tflog.Info(ctx, fmt.Sprintf("Private network %s not found, removing from state", firewallPolicyId))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading private network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving private network. Please try again or report this issue to the provider developers",
		)
		return
	}

	readModel, diags := models.NewPrivateNetworkServerAssignModel(ctx, *apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *PrivateNetworkServersResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource does not support updates. Changes will trigger resource replacement.",
	)
}

func (r *PrivateNetworkServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PrivateNetworkServerAssignModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	privateNetworkId := data.Id.ValueString()

	var serverIds []string
	diags := data.Servers.ElementsAs(ctx, &serverIds, false)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if len(serverIds) == 0 {
		tflog.Info(ctx, fmt.Sprintf("No servers to remove from private network %s; removing from state only", privateNetworkId))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Removing %d servers from private network %s", len(serverIds), privateNetworkId))

	for _, serverId := range serverIds {
		tflog.Info(ctx, fmt.Sprintf("Removing server %s from private network %s", serverId, privateNetworkId))

		_, err := r.client.DeletePrivateNetworkServer(privateNetworkId, serverId)
		if err != nil {
			if errors.Is(err, util.ErrNotFound) {
				tflog.Info(ctx, fmt.Sprintf("Server %s or private network %s not found - may already be removed", serverId, privateNetworkId))
				continue
			}

			resp.Diagnostics.AddError(
				"Error removing server from private network",
				fmt.Sprintf("Private network %s, server %s: %s", privateNetworkId, serverId, err.Error()),
			)
			return
		}

		timeouts := util.GetResourceTimeouts("PRIVATE_NETWORKS_OPERATIONS")
		waitOptions := util.NewWaitOptions(
			timeouts.Default,
			timeouts.RetryInterval,
			timeouts.MinTimeout,
			[]string{util.StatusConfiguring},
			[]string{util.StateActive},
		)

		_, waitDiags := util.WaitForResourceState(ctx, privateNetworkId, r.client, waitOptions)
		resp.Diagnostics.Append(waitDiags...)
		if resp.Diagnostics.HasError() {
			tflog.Error(ctx, "Wait for private network state after remove failed")
			return
		}
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully removed %d servers from private network %s", len(serverIds), privateNetworkId))
}
