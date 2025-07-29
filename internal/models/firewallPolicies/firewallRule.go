package firewallPolicies

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallRuleResponse struct {
	ID          string  `json:"id"`
	Protocol    string  `json:"protocol"`
	PortFrom    int     `json:"port_from"`
	PortTo      int     `json:"port_to"`
	Source      string  `json:"source"`
	Description *string `json:"description"`
	Action      string  `json:"action"`
}

func NewFirewallRuleObject(rule FirewallRuleResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	var description types.String
	if rule.Description != nil {
		description = types.StringValue(*rule.Description)
	} else {
		description = types.StringNull()
	}

	attrs := map[string]attr.Value{
		"id":          types.StringValue(rule.ID),
		"protocol":    types.StringValue(rule.Protocol),
		"port_from":   types.Int64Value(int64(rule.PortFrom)),
		"port_to":     types.Int64Value(int64(rule.PortTo)),
		"source":      types.StringValue(rule.Source),
		"description": description,
		"action":      types.StringValue(rule.Action),
	}

	obj, objDiags := types.ObjectValue(FirewallRuleObjectType().AttrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}

func NewFirewallRulesList(rules []FirewallRuleResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(rules) == 0 {
		return types.ListValueMust(FirewallRuleObjectType(), []attr.Value{}), diags
	}

	elements := make([]attr.Value, 0, len(rules))

	for _, rule := range rules {
		ruleObj, objDiags := NewFirewallRuleObject(rule)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			elements = append(elements, ruleObj)
		}
	}

	if diags.HasError() {
		return types.ListValueMust(FirewallRuleObjectType(), []attr.Value{}), diags
	}

	list, listDiags := types.ListValue(FirewallRuleObjectType(), elements)
	diags.Append(listDiags...)

	return list, diags
}

func FirewallRuleObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"protocol":    types.StringType,
			"port_from":   types.Int64Type,
			"port_to":     types.Int64Type,
			"source":      types.StringType,
			"description": types.StringType,
			"action":      types.StringType,
		},
	}
}

func FirewallRuleDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Rule identifier",
		},
		"protocol": schema.StringAttribute{
			Computed:    true,
			Description: "Internet protocol (TCP, UDP, ICMP, TCP/UDP, IPSEC, GRE)",
		},
		"port_from": schema.Int64Attribute{
			Computed:    true,
			Description: "First port in range (1-65535)",
		},
		"port_to": schema.Int64Attribute{
			Computed:    true,
			Description: "Second port in range (1-65535)",
		},
		"source": schema.StringAttribute{
			Computed:    true,
			Description: "Source IP address or range",
		},
		"description": schema.StringAttribute{
			Computed:    true,
			Description: "Rule description",
		},
		"action": schema.StringAttribute{
			Computed:    true,
			Description: "Rule action (allow/deny)",
		},
	}
}

func FirewallRuleResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed:    true,
			Description: "Rule identifier",
		},
		"protocol": rschema.StringAttribute{
			Computed:    true,
			Description: "Internet protocol (TCP, UDP, ICMP, TCP/UDP, IPSEC, GRE)",
		},
		"port_from": rschema.Int64Attribute{
			Computed:    true,
			Description: "First port in range (1-65535)",
		},
		"port_to": rschema.Int64Attribute{
			Computed:    true,
			Description: "Second port in range (1-65535)",
		},
		"source": rschema.StringAttribute{
			Computed:    true,
			Description: "Source IP address or range",
		},
		"description": rschema.StringAttribute{
			Computed:    true,
			Description: "Rule description",
		},
		"action": rschema.StringAttribute{
			Computed:    true,
			Description: "Rule action (allow/deny)",
		},
	}
}
