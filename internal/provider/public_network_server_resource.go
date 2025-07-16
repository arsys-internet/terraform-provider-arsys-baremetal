package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicNetwork"
	"terraform-provider-arsys-baremetal/internal/util"
)

var _ resource.Resource = &PublicNetworkServerResource{}

func NewPublicNetworkServerResource() resource.Resource {
	return &PublicNetworkServerResource{}
}

type PublicNetworkServerResource struct {
	client service.ApiPublicNetworkServiceInterface
}

func (r *PublicNetworkServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network_server"
}

func (r *PublicNetworkServerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PublicNetworkServerSchema(ctx)
}

func (r *PublicNetworkServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *PublicNetworkServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convertir los tipos de Terraform a strings
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
			fmt.Sprintf("Could not assign servers to public network: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getPublicNetworkTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
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

func (r *PublicNetworkServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Verificar que la red pública existe
	publicNetwork, err := r.client.GetPublicNetwork(data.PublicNetworkId.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading public network",
			fmt.Sprintf("Could not read public network: %s", err),
		)
		return
	}

	if publicNetwork == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *PublicNetworkServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
			fmt.Sprintf("Could not update servers in public network: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getPublicNetworkTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
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

func (r *PublicNetworkServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PublicNetworkServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Para eliminar, enviamos una lista vacía de servidores
	request := &models.PublicNetworkServerRequest{
		Servers: []string{},
	}

	err := r.client.AssignServersToPublicNetwork(data.PublicNetworkId.ValueString(), request)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing servers from public network",
			fmt.Sprintf("Could not remove servers from public network: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getPublicNetworkTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
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
