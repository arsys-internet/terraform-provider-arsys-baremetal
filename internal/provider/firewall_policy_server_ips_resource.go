package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewall_policy"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &FirewallPolicyServerIPsAssignResource{}

func NewFirewallPolicyServerIPsAssignResource() resource.Resource {
	return &FirewallPolicyServerIPsAssignResource{}
}

type FirewallPolicyServerIPsAssignResource struct {
	client *service.ApiFirewallPolicyService
}

func (r *FirewallPolicyServerIPsAssignResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy_server_ips"
}

func (r *FirewallPolicyServerIPsAssignResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyAssignmentResourceSchema(ctx)
}

func (r *FirewallPolicyServerIPsAssignResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetFirewallPolicyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	policyService, ok := client.(*service.ApiFirewallPolicyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = policyService
}

func (r *FirewallPolicyServerIPsAssignResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.FirewallPolicyServerModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyID := data.Id.ValueString()

	assignRequest, diags := data.ToAssignRequest(ctx)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	apiResponse, assignErr := r.client.AssignServerIPsToFirewallPolicy(firewallPolicyID, assignRequest)
	if assignErr != nil {
		resp.Diagnostics.AddError(
			"Error assigning server IPs to firewall policy",
			fmt.Sprintf("Error: %s", assignErr.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("FIREWALL_POLICY_OPERATIONS")
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
		tflog.Error(ctx, "Wait for firewall policy state failed")
		return
	}

	finalPolicy, fwErr := r.client.GetFirewallPolicy(id)
	if fwErr != nil {
		resp.Diagnostics.AddError(
			"Error getting final firewall policy state",
			fmt.Sprintf("Could not get final firewall policy state: %s", fwErr.Error()),
		)
		return
	}

	finalModel, diags := models.NewFirewallPolicyServerModel(ctx, *finalPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully assigned %d server IPs to firewall policy: %s", len(assignRequest.ServerIPs), firewallPolicyID))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *FirewallPolicyServerIPsAssignResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.FirewallPolicyServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyId := data.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Reading firewall policy assignment: %s", firewallPolicyId))

	apiResponse, err := r.client.GetFirewallPolicy(firewallPolicyId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Firewall policy %s not found, removing from state", firewallPolicyId))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Firewall policy %s not found, removing from state", firewallPolicyId))
		resp.State.RemoveResource(ctx)
		return
	}

	readModel, diags := models.NewFirewallPolicyServerModel(ctx, *apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *FirewallPolicyServerIPsAssignResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource uses RequiresReplace for all changes. Any modifications should result in destroy + create, not update. Please check your Terraform configuration.",
	)
}

func (r *FirewallPolicyServerIPsAssignResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.FirewallPolicyServerModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyID := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Removing firewall policy assignment from Terraform state only: %s", firewallPolicyID))
	tflog.Info(ctx, "Note: Server IPs remain assigned to the firewall policy - only removing from Terraform management")
}
