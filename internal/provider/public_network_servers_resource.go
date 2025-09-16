package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/public_network"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	waitResult, diags := util.WaitForResourceState(
		ctx,
		data.PublicNetworkId.ValueString(),
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if waitResult == nil {
		resp.Diagnostics.AddError(
			"Resource state timeout",
			"Public network did not reach active state within timeout period",
		)
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s-%s", data.PublicNetworkId.ValueString(), strings.Join(servers, "-")))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	//TODO: refactor this to return the public network
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

	waitResult, diags := util.WaitForResourceState(
		ctx,
		data.PublicNetworkId.ValueString(),
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if waitResult == nil {
		resp.Diagnostics.AddError(
			"Resource state timeout",
			"Public network did not reach active state within timeout period",
		)
		return
	}

	data.Id = types.StringValue(fmt.Sprintf("%s-%s", data.PublicNetworkId.ValueString(), strings.Join(servers, "-")))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
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

	waitResult, diags := util.WaitForResourceState(
		ctx,
		data.PublicNetworkId.ValueString(),
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if waitResult == nil {
		resp.Diagnostics.AddError(
			"Resource state timeout",
			"Public network did not reach active state within timeout period",
		)
		return
	}
}
