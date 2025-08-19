package firewallPolicies

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallRuleResponse struct {
	Id          string  `json:"id"`
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
		"id":          types.StringValue(rule.Id),
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

type FirewallRuleCreateRequest struct {
	Protocol    string  `json:"protocol"`
	PortFrom    int     `json:"port_from"`
	PortTo      int     `json:"port_to"`
	Source      string  `json:"source,omitempty"`
	Description *string `json:"description,omitempty"`
	Action      *string `json:"action,omitempty"`
}

func ConvertRulesToCreateRequest(rules types.List) ([]FirewallRuleCreateRequest, error) {
	if rules.IsNull() || rules.IsUnknown() {
		return []FirewallRuleCreateRequest{}, nil
	}

	var rulesObjects []types.Object
	if diags := rules.ElementsAs(context.Background(), &rulesObjects, false); diags.HasError() {
		return nil, fmt.Errorf("failed to convert rules: %v", diags.Errors())
	}

	result := make([]FirewallRuleCreateRequest, 0, len(rulesObjects))

	for i, ruleObj := range rulesObjects {
		attrs := ruleObj.Attributes()
		rule := FirewallRuleCreateRequest{}

		if protocolVal, ok := attrs["protocol"].(types.String); !ok || protocolVal.IsNull() || protocolVal.ValueString() == "" {
			return nil, fmt.Errorf("rule[%d]: protocol is required", i)
		} else {
			rule.Protocol = protocolVal.ValueString()
		}

		if portFromVal, ok := attrs["port_from"].(types.Int64); !ok || portFromVal.IsNull() {
			return nil, fmt.Errorf("rule[%d]: port_from is required", i)
		} else {
			rule.PortFrom = int(portFromVal.ValueInt64())
		}

		if portToVal, ok := attrs["port_to"].(types.Int64); !ok || portToVal.IsNull() {
			return nil, fmt.Errorf("rule[%d]: port_to is required", i)
		} else {
			rule.PortTo = int(portToVal.ValueInt64())
		}

		if sourceVal, ok := attrs["source"].(types.String); ok && !sourceVal.IsNull() {
			rule.Source = sourceVal.ValueString()
		}

		if descVal, ok := attrs["description"].(types.String); ok && !descVal.IsNull() && descVal.ValueString() != "" {
			desc := descVal.ValueString()
			rule.Description = &desc
		}

		if actionVal, ok := attrs["action"].(types.String); ok && !actionVal.IsNull() && actionVal.ValueString() != "" {
			action := actionVal.ValueString()
			rule.Action = &action
		}

		result = append(result, rule)
	}

	return result, nil
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
		"protocol": rschema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Internet protocol",
			Validators: []validator.String{
				stringvalidator.OneOf("TCP", "UDP", "ICMP", "AH", "ESP", "GRE"),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"port_from": rschema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "First port in range. Required for UDP and TCP protocols",
			Validators: []validator.Int64{
				int64validator.Between(1, 65535),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"port_to": rschema.Int64Attribute{
			Optional:    true,
			Computed:    true,
			Description: "Second port in range. Required for UDP and TCP protocols",
			Validators: []validator.Int64{
				int64validator.AtMost(65535),
			},
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
		},
		"source": rschema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Default:     stringdefault.StaticString("0.0.0.0"),
			Description: "IPs from which access is available. Setting 0.0.0.0 all IPs are allowed",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"description": rschema.StringAttribute{
			Optional:    true,
			Computed:    true,
			Description: "Rule description",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"action": rschema.StringAttribute{
			Computed:    true,
			Optional:    true,
			Description: "Rule action (allow/deny)",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"id": rschema.StringAttribute{
			Computed:    true,
			Description: "Rule identifier",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
	}
}

func ValidateFirewallRules(rules types.List, basePath path.Path) diag.Diagnostics {
	var diags diag.Diagnostics

	if rules.IsNull() || rules.IsUnknown() {
		return diags
	}

	var rulesObjects []types.Object
	if ruleDiags := rules.ElementsAs(context.Background(), &rulesObjects, false); ruleDiags.HasError() {
		diags.Append(ruleDiags...)
		return diags
	}

	for i, ruleObj := range rulesObjects {
		attrs := ruleObj.Attributes()
		ruleBasePath := basePath.AtListIndex(i)

		if protocolVal, ok := attrs["protocol"].(types.String); !ok || protocolVal.IsNull() || protocolVal.ValueString() == "" {
			diags.AddAttributeError(
				ruleBasePath.AtName("protocol"),
				"Missing required field",
				fmt.Sprintf("'protocol' field is required in rule[%d]", i),
			)
		}

		if portFromVal, ok := attrs["port_from"].(types.Int64); !ok || portFromVal.IsNull() {
			diags.AddAttributeError(
				ruleBasePath.AtName("port_from"),
				"Missing required field",
				fmt.Sprintf("'port_from' field is required in rule[%d]", i),
			)
		}

		if portToVal, ok := attrs["port_to"].(types.Int64); !ok || portToVal.IsNull() {
			diags.AddAttributeError(
				ruleBasePath.AtName("port_to"),
				"Missing required field",
				fmt.Sprintf("'port_to' field is required in rule[%d]", i),
			)
		}
	}

	return diags
}
