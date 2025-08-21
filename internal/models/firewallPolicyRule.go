package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/firewallPolicies"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPolicyRuleModel struct {
	FirewallPolicyId types.String `tfsdk:"firewall_policy_id"`
	Id               types.String `tfsdk:"id"`
	Protocol         types.String `tfsdk:"protocol"`
	PortFrom         types.Int64  `tfsdk:"port_from"`
	PortTo           types.Int64  `tfsdk:"port_to"`
	Source           types.String `tfsdk:"source"`
	Description      types.String `tfsdk:"description"`
	Action           types.String `tfsdk:"action"`
}

func NewFirewallPolicyRuleModel(_ context.Context, firewallPolicyId string, rule firewallPolicies.FirewallRuleResponse) *FirewallPolicyRuleModel {
	var description types.String
	if rule.Description != nil {
		description = types.StringValue(*rule.Description)
	} else {
		description = types.StringNull()
	}

	model := &FirewallPolicyRuleModel{
		Id:               types.StringValue(rule.Id),
		FirewallPolicyId: types.StringValue(firewallPolicyId),
		Protocol:         types.StringValue(rule.Protocol),
		PortFrom:         types.Int64Value(int64(rule.PortFrom)),
		PortTo:           types.Int64Value(int64(rule.PortTo)),
		Source:           types.StringValue(rule.Source),
		Description:      description,
		Action:           types.StringValue(rule.Action),
	}

	return model
}

func FirewallPolicyRuleDataSourceSchema(_ context.Context) schema.Schema {
	baseAttributes := firewallPolicies.FirewallRuleDataSourceSchema()

	baseAttributes["id"] = schema.StringAttribute{
		Required:    true,
		Description: "Rule identifier",
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(util.HexID32Pattern),
				"must be a valid Id",
			),
		},
	}

	baseAttributes["firewall_policy_id"] = schema.StringAttribute{
		Required:    true,
		Description: "Firewall policy identifier",
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(util.HexID32Pattern),
				"must be a valid Id",
			),
		},
	}

	return schema.Schema{
		Description: "Fetches information about a specific firewall rule",
		Attributes:  baseAttributes,
	}
}
