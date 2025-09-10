package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewallPolicy"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &FirewallPolicyRemoveRuleResource{}

func NewFirewallPolicyRemoveRuleResource() resource.Resource {
	return &FirewallPolicyRemoveRuleResource{}
}

type FirewallPolicyRemoveRuleResource struct {
	client *service.ApiFirewallPolicyService
}

func (r *FirewallPolicyRemoveRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy_rule_remove"
}

func (r *FirewallPolicyRemoveRuleResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyRuleRemoveResourceSchema(ctx)
}

func (r *FirewallPolicyRemoveRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FirewallPolicyRemoveRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.FirewallPolicyRuleRemoveResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyId := data.Id.ValueString()
	ruleId := data.RuleId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Removing rule %s from firewall policy %s", ruleId, firewallPolicyId))

	_, err := r.client.DeleteFirewallPolicyRule(firewallPolicyId, ruleId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Rule %s or policy %s not found - may already be deleted", ruleId, firewallPolicyId))
		} else {
			resp.Diagnostics.AddError(
				"Error removing rule from firewall policy",
				fmt.Sprintf("Error: %s", err.Error()),
			)
			return
		}
	}

	timeouts := util.GetResourceTimeouts("FIREWALL_POLICY_OPERATIONS")
	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StateActive},
	)

	_, diags := util.WaitForResourceState(ctx, firewallPolicyId, r.client, waitOptions)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Wait for firewall policy state failed")
		return
	}

	finalPolicy, err := r.client.GetFirewallPolicy(firewallPolicyId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting final firewall policy state",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if finalPolicy == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving firewall policy after rule removal. Please report this issue to the provider developers.",
		)
		return
	}

	finalModel, diags := models.NewFirewallPolicyRuleRemoveResourceModel(ctx, ruleId, *finalPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully removed rule %s from firewall policy %s", ruleId, firewallPolicyId))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *FirewallPolicyRemoveRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.FirewallPolicyRuleRemoveResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicy, err := r.client.GetFirewallPolicy(data.Id.ValueString())
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if firewallPolicy == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	updatedModel, diags := models.NewFirewallPolicyRuleRemoveResourceModel(
		ctx,
		data.RuleId.ValueString(),
		*firewallPolicy,
	)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, updatedModel)...)
}

func (r *FirewallPolicyRemoveRuleResource) Update(_ context.Context, _ resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError(
		"Update not supported",
		"This resource does not support updates. Changes will trigger resource replacement.",
	)
}

func (r *FirewallPolicyRemoveRuleResource) Delete(ctx context.Context, _ resource.DeleteRequest, _ *resource.DeleteResponse) {
	tflog.Info(ctx, "Removing rule removal tracking from Terraform state")
}
