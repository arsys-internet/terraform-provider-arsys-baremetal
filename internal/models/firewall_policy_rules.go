package models

import (
	"context"
	"regexp"
	firewallpolicy "terraform-provider-arsys-baremetal/internal/models/firewall_policy"

	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPolicyRulesModel struct {
	Id    types.String `tfsdk:"id"`
	Rules types.List   `tfsdk:"rules"`
}

func NewFirewallPolicyRulesModel(_ context.Context, policyID string, rules []firewallpolicy.FirewallRuleResponse) (*FirewallPolicyRulesModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	rulesList, rulesDiags := firewallpolicy.NewFirewallRulesList(rules)
	diags.Append(rulesDiags...)

	model := &FirewallPolicyRulesModel{
		Id:    types.StringValue(policyID),
		Rules: rulesList,
	}

	return model, diags
}

func FirewallPolicyRulesSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Firewall policy rules",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Id of the firewall policy",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid firewall policy ID",
					),
				},
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Rules assigned to firewall policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: firewallpolicy.FirewallRuleDataSourceSchema(),
				},
			},
		},
	}
}
