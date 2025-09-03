package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicNetwork"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &PublicNetworkIpResource{}
	_ resource.ResourceWithConfigure   = &PublicNetworkIpResource{}
	_ resource.ResourceWithImportState = &PublicNetworkIpResource{}
)

func NewPublicNetworkIpResource() resource.Resource {
	return &PublicNetworkIpResource{}
}

type PublicNetworkIpResource struct {
	client *service.ApiPublicNetworkService
}

func (r *PublicNetworkIpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network_ip"
}

func (r *PublicNetworkIpResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.PublicNetworkIpResourceSchema(ctx)
}

func (r *PublicNetworkIpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

	PublicNetworkIpService, ok := client.(*service.ApiPublicNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = PublicNetworkIpService
}

func (r *PublicNetworkIpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.PublicNetworkIpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	ips := data.Ips

	publicNetworkId := data.PublicNetworkId.ValueString()

	for _, ipId := range ips {
		id := ipId.ValueString()

		tflog.Info(ctx, fmt.Sprintf("Reading IP %s in the public network with ID: %s", id, publicNetworkId))

		apiResponse, err := r.client.GetPublicNetworkIp(publicNetworkId, id)
		if err != nil {
			if strings.Contains(err.Error(), "not found") {
				tflog.Info(ctx, fmt.Sprintf("IP %s not found in the Public network with ID %s, removing from state", id, publicNetworkId))
				resp.State.RemoveResource(ctx)
				return
			}

			resp.Diagnostics.AddError(
				"Error reading public network",
				fmt.Sprintf("Error: %s", err.Error()),
			)
			return
		}

		if apiResponse == nil {
			tflog.Info(ctx, fmt.Sprintf("IP %s not found in the public network %s not found, removing from state", id, publicNetworkId))
			resp.State.RemoveResource(ctx)
			return
		}

		readModel, diags := models.NewPublicNetworkIpModel(ctx, apiResponse)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
	}
}

func (r *PublicNetworkIpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.PublicNetworkIpResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.PublicNetworkId.IsNull() || data.PublicNetworkId.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("datacenter_id"),
			"Missing required field",
			"Either 'datacenter_id' field is required when associating an IP to public network",
		)
	}

	if data.Action.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("action"),
			"Missing required field",
			"Either 'action' field is required when associating an IP to public network",
		)
	}

	if len(data.Ips) == 0 {
		resp.Diagnostics.AddAttributeError(
			path.Root("ips"),
			"Missing required field",
			"Either 'ips' field is required when associating an IP to public network",
		)
	}

	createRequest := data.ToCreateRequest()

	apiResponse, err := r.client.AssignIpToPublicNetwork(data.PublicNetworkId.ValueString(), &createRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error associating an IP to public network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	model, diags := models.NewPublicNetworkIpResourceModel(ctx, &data, *apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//timeouts := util.GetResourceTimeouts("PUBLIC_NETWORK")
	//
	//waitOptions := util.NewWaitOptions(
	//	timeouts.Default,
	//	timeouts.RetryInterval,
	//	timeouts.MinTimeout,
	//	[]string{util.StateDeploying},
	//	[]string{util.StatePoweredOn, util.StatePoweredOff},
	//)
	//
	//waitResult, diags := util.WaitForResourceState(
	//	ctx,
	//	apiResponse.Id,
	//	r.client,
	//	waitOptions,
	//)
	//
	//resp.Diagnostics.Append(diags...)
	//if resp.Diagnostics.HasError() {
	//	return
	//}
	//
	//finalModel := model
	//if waitResult != nil && waitResult.Resource != nil {
	//	if PublicNetworkIpModel, ok := waitResult.Resource.(*models.PublicNetworkIpModel); ok {
	//		finalModel = PublicNetworkIpModel
	//		tflog.Info(ctx, fmt.Sprintf("Public network reached final state: %s", waitResult.FinalState))
	//	}
	//}

	//resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)

	data.Id = types.StringValue(fmt.Sprintf("%s-%s", data.PublicNetworkId.ValueString(), strings.Join(createRequest.Ips, "-")))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}

func (r *PublicNetworkIpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	createReq := resource.CreateRequest{
		Plan: req.Plan,
	}

	createResp := &resource.CreateResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Create(ctx, createReq, createResp)

	resp.State = createResp.State
	resp.Diagnostics = createResp.Diagnostics
}

func (r *PublicNetworkIpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.PublicNetworkIpResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetworkId := data.PublicNetworkId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Removing public network IP assignment from Terraform state only: %s", publicNetworkId))
	tflog.Info(ctx, "Note: Server IPs remain assigned to the public network - only removing from Terraform management")
}

func (r *PublicNetworkIpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
