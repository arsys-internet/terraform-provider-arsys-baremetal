package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	"terraform-provider-arsys-baremetal/internal/models/firewallPolicies"
	service "terraform-provider-arsys-baremetal/internal/services/firewallPolicy"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &FirewallPolicyRuleResource{}

func NewFirewallPolicyRuleResource() resource.Resource {
	return &FirewallPolicyRuleResource{}
}

type FirewallPolicyRuleResource struct {
	client *service.ApiFirewallPolicyService
}

func (r *FirewallPolicyRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy_rule_add"
}

func (r *FirewallPolicyRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyRuleAddResourceSchema(ctx)
}

func (r *FirewallPolicyRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallPolicyRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.FirewallPolicyRuleAddResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	if !data.Rules.IsNull() {
		rulesDiags := firewallPolicies.ValidateFirewallRules(data.Rules, path.Root("rules"))
		resp.Diagnostics.Append(rulesDiags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	createRequest, err := data.ToAddRequest(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error converting model to request",
			fmt.Sprintf("Could not convert rules to request: %s", err.Error()),
		)
		return
	}

	apiResponse, assignErr := r.client.CreateFirewallPolicyRule(id, createRequest)
	if assignErr != nil {
		resp.Diagnostics.AddError(
			"Error adding new rules to firewall policy",
			fmt.Sprintf("Error: %s", assignErr),
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

	apiResponse.Id = id

	_, diags := util.WaitForResourceState(
		ctx,
		id,
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Wait for firewall policy state failed")
		return
	}

	finalPolicy, fwErr := r.client.GetFirewallPolicy(id)
	if fwErr != nil {
		resp.Diagnostics.AddError(
			"Error getting final firewall policy state",
			fmt.Sprintf("Error: %s", fwErr.Error()),
		)
		return
	}

	finalModel, diags := models.NewFirewallPolicyRuleResourceModel(ctx, data.Rules, *finalPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully added rules %d to firewall policy: %s", len(createRequest.Rules), id))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *FirewallPolicyRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.FirewallPolicyRuleAddResourceModel

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

	readModel, diags := models.NewFirewallPolicyRuleResourceModel(ctx, data.Rules, *apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *FirewallPolicyRuleResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource cannot be updated. Please check your Terraform configuration.",
	)
}

func (r *FirewallPolicyRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.FirewallPolicyRuleAddResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyID := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Removing firewall policy assignment from Terraform state only: %s", firewallPolicyID))
	tflog.Info(ctx, "Note: Server IPs remain assigned to the firewall policy - only removing from Terraform management")
}
