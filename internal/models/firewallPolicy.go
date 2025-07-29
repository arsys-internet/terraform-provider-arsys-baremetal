package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"terraform-provider-arsys-baremetal/internal/models/firewallPolicies"
)

type FirewallPolicyResponse struct {
	ID           string                                      `json:"id"`
	Name         string                                      `json:"name"`
	Description  string                                      `json:"description"`
	State        string                                      `json:"state"`
	CreationDate string                                      `json:"creation_date"`
	Default      int                                         `json:"default"`
	Rules        []firewallPolicies.FirewallRuleResponse     `json:"rules"`
	ServerIPs    []firewallPolicies.FirewallServerIPResponse `json:"server_ips"`
	CloudPanelID string                                      `json:"cloudpanel_id"`
}

type FirewallPolicyModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	State        types.String `tfsdk:"state"`
	CreationDate types.String `tfsdk:"creation_date"`
	Default      types.Int64  `tfsdk:"default"`
	Rules        types.List   `tfsdk:"rules"`
	ServerIPs    types.List   `tfsdk:"server_ips"`
	CloudPanelID types.String `tfsdk:"cloudpanel_id"`
}

func NewFirewallPolicyModel(_ context.Context, fp FirewallPolicyResponse) (*FirewallPolicyModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	rulesList, rulesDiags := firewallPolicies.NewFirewallRulesList(fp.Rules)
	diags.Append(rulesDiags...)

	serverIPsList, serverIPsDiags := firewallPolicies.NewFirewallServerIPsList(fp.ServerIPs)
	diags.Append(serverIPsDiags...)

	model := &FirewallPolicyModel{
		ID:           types.StringValue(fp.ID),
		Name:         types.StringValue(fp.Name),
		Description:  types.StringValue(fp.Description),
		State:        types.StringValue(fp.State),
		CreationDate: types.StringValue(fp.CreationDate),
		Default:      types.Int64Value(int64(fp.Default)),
		Rules:        rulesList,
		ServerIPs:    serverIPsList,
		CloudPanelID: types.StringValue(fp.CloudPanelID),
	}

	return model, diags
}

func FirewallPolicyObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":            types.StringType,
			"name":          types.StringType,
			"description":   types.StringType,
			"state":         types.StringType,
			"creation_date": types.StringType,
			"default":       types.Int64Type,
			"rules":         types.ListType{ElemType: firewallPolicies.FirewallRuleObjectType()},
			"server_ips":    types.ListType{ElemType: firewallPolicies.FirewallServerIPObjectType()},
			"cloudpanel_id": types.StringType,
		},
	}
}

func FirewallPolicyDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Fetches information about a specific firewall policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall policy identifier",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy description",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy state",
			},
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when firewall policy was created",
			},
			"default": schema.Int64Attribute{
				Computed:    true,
				Description: "Define default panel firewalls (1 = default, 0 = custom)",
			},
			"cloudpanel_id": schema.StringAttribute{
				Computed:    true,
				Description: "Public identifier shown in panel",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Firewall policy rules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: firewallPolicies.FirewallRuleDataSourceSchema(),
				},
			},
			"server_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Servers assigned to firewall policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: firewallPolicies.FirewallServerIPDataSourceSchema(),
				},
			},
		},
	}
}
