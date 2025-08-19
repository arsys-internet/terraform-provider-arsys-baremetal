package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/firewallPolicies"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPolicyRuleRemoveResourceModel struct {
	RuleId       types.String `tfsdk:"rule_id"`
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	State        types.String `tfsdk:"state"`
	CreationDate types.String `tfsdk:"creation_date"`
	Default      types.Int64  `tfsdk:"default"`
	Rules        types.List   `tfsdk:"rules"`
	ServerIPs    types.List   `tfsdk:"server_ips"`
	CloudPanelID types.String `tfsdk:"cloudpanel_id"`
}

func NewFirewallPolicyRuleRemoveResourceModel(_ context.Context, ruleId string, fp FirewallPolicyResponse) (*FirewallPolicyRuleRemoveResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	rulesList, rulesDiags := firewallPolicies.NewFirewallRulesList(fp.Rules)
	diags.Append(rulesDiags...)

	serverIPsList, serverIPsDiags := firewallPolicies.NewFirewallServerIPsList(fp.ServerIPs)
	diags.Append(serverIPsDiags...)

	var description types.String
	if fp.Description != nil {
		description = types.StringValue(*fp.Description)
	} else {
		description = types.StringNull()
	}

	model := &FirewallPolicyRuleRemoveResourceModel{
		RuleId:       types.StringValue(ruleId),
		Id:           types.StringValue(fp.Id),
		Name:         types.StringValue(fp.Name),
		Description:  description,
		State:        types.StringValue(fp.State),
		CreationDate: types.StringValue(fp.CreationDate),
		Default:      types.Int64Value(int64(fp.Default)),
		Rules:        rulesList,
		ServerIPs:    serverIPsList,
		CloudPanelID: types.StringValue(fp.CloudPanelID),
	}

	return model, diags
}

func FirewallPolicyRuleRemoveResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Assigns rules to an existing firewall policy",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Required:    true,
				Description: "Resource identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid Id",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"rule_id": rschema.StringAttribute{
				Required:    true,
				Description: "Rule identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid Id",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},

			"name": rschema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": rschema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy description",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy state",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "Date when firewall policy was created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default": rschema.Int64Attribute{
				Computed:    true,
				Description: "Define default panel firewalls",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the cloud panel",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"rules": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "Complete list of rules in the firewall policy after assignment",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: rschema.NestedAttributeObject{
					Attributes: firewallPolicies.FirewallRuleResourceSchema(),
				},
			},
			"server_ips": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "Servers assigned to firewall policy",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: rschema.NestedAttributeObject{
					Attributes: firewallPolicies.FirewallServerIPResourceSchema(),
				},
			},
		},
	}
}
