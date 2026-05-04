package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicnetwork"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var _ resource.Resource = &PublicNetworkServersResource{}

func NewPublicNetworkServerResource() resource.Resource {
	return &PublicNetworkServersResource{}
}

type PublicNetworkServersResource struct {
	client service.ApiPublicNetworkServiceInterface
}

func (r *PublicNetworkServersResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network_servers"
}

func (r *PublicNetworkServersResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PublicNetworkServerSchema(ctx)
}

func (r *PublicNetworkServersResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	r.client = client
}

func (r *PublicNetworkServersResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	servers := make([]string, len(data.Servers))
	for i, server := range data.Servers {
		servers[i] = server.ValueString()
	}

	request := &models.PublicNetworkServerRequest{
		Servers: servers,
	}

	err := r.client.AssignServersToPublicNetwork(data.PublicNetworkId.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error assigning servers to public network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("PUBLIC_NETWORK")

	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StatePoweredOn, util.StatePoweredOff},
	)

	_, diags := util.WaitForResourceState(
		ctx,
		data.PublicNetworkId.ValueString(),
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	finalPublicNetwork, fwErr := r.client.GetPublicNetwork(data.PublicNetworkId.ValueString())
	if fwErr != nil {
		resp.Diagnostics.AddError(
			"Error getting final public network state",
			fmt.Sprintf("Error: %s", fwErr.Error()),
		)
		return
	}

	if finalPublicNetwork == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving public network after assign servers. Please report this issue to the provider developers.",
		)
		return
	}

	finalModel, diags := models.NewPublicNetworkServerResourceModel(
		ctx,
		data.PublicNetworkId.ValueString(),
		servers,
		finalPublicNetwork,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *PublicNetworkServersResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetwork, err := r.client.GetPublicNetwork(data.PublicNetworkId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading public network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if publicNetwork == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	servers := make([]string, len(publicNetwork.Servers))
	for i, server := range publicNetwork.Servers {
		servers[i] = server.Id
	}

	finalModel, diags := models.NewPublicNetworkServerResourceModel(
		ctx,
		data.PublicNetworkId.ValueString(),
		servers,
		publicNetwork,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *PublicNetworkServersResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	servers := make([]string, len(data.Servers))
	for i, server := range data.Servers {
		servers[i] = server.ValueString()
	}

	request := &models.PublicNetworkServerRequest{
		Servers: servers,
	}

	err := r.client.AssignServersToPublicNetwork(data.PublicNetworkId.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating servers in public network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("PUBLIC_NETWORK")

	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StatePoweredOn, util.StatePoweredOff},
	)

	_, diags := util.WaitForResourceState(
		ctx,
		data.PublicNetworkId.ValueString(),
		r.client,
		waitOptions,
	)

	finalPublicNetwork, fwErr := r.client.GetPublicNetwork(data.PublicNetworkId.ValueString())
	if fwErr != nil {
		resp.Diagnostics.AddError(
			"Error getting final public network state",
			fmt.Sprintf("Error: %s", fwErr.Error()),
		)
		return
	}

	if finalPublicNetwork == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving public network after assign servers. Please report this issue to the provider developers.",
		)
		return
	}

	finalModel, diags := models.NewPublicNetworkServerResourceModel(
		ctx,
		data.PublicNetworkId.ValueString(),
		servers,
		finalPublicNetwork,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *PublicNetworkServersResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	request := &models.PublicNetworkServerRequest{
		Servers: []string{},
	}

	err := r.client.AssignServersToPublicNetwork(data.PublicNetworkId.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing servers from public network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("PUBLIC_NETWORK")

	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StatePoweredOn, util.StatePoweredOff},
	)

	_, diags := util.WaitForResourceState(
		ctx,
		data.PublicNetworkId.ValueString(),
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	finalPublicNetwork, fwErr := r.client.GetPublicNetwork(data.PublicNetworkId.ValueString())
	if fwErr != nil {
		resp.Diagnostics.AddError(
			"Error getting final public network state",
			fmt.Sprintf("Error: %s", fwErr.Error()),
		)
		return
	}

	if finalPublicNetwork == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving public network after assign servers. Please report this issue to the provider developers.",
		)
		return
	}

	var servers []string
	finalModel, diags := models.NewPublicNetworkServerResourceModel(
		ctx,
		data.PublicNetworkId.ValueString(),
		servers,
		finalPublicNetwork,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}
