package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/firewallPolicies"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPolicyServerAssignRequest struct {
	ServerIPs []string `json:"server_ips"`
}

type FirewallPolicyServerModel struct {
	Id              types.String `tfsdk:"id"`
	ServerIPs       types.Set    `tfsdk:"server_ips"`
	Name            types.String `tfsdk:"name"`
	Description     types.String `tfsdk:"description"`
	State           types.String `tfsdk:"state"`
	CreationDate    types.String `tfsdk:"creation_date"`
	Default         types.Int64  `tfsdk:"default"`
	Rules           types.List   `tfsdk:"rules"`
	ServerIPsDetail types.List   `tfsdk:"server_ips_detail"`
	CloudPanelId    types.String `tfsdk:"cloudpanel_id"`
}

func NewFirewallPolicyServerModel(_ context.Context, fp FirewallPolicyResponse) (*FirewallPolicyServerModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	rulesList, rulesDiags := firewallPolicies.NewFirewallRulesList(fp.Rules)
	diags.Append(rulesDiags...)

	serverIPIds := make([]string, len(fp.ServerIPs))
	for i, serverIP := range fp.ServerIPs {
		serverIPIds[i] = serverIP.Id
	}
	serverIPsSet, _ := types.SetValueFrom(context.Background(), types.StringType, serverIPIds)

	serverIPsList, serverIPsDiags := firewallPolicies.NewFirewallServerIPsList(fp.ServerIPs)
	diags.Append(serverIPsDiags...)

	var description types.String
	if fp.Description != nil {
		description = types.StringValue(*fp.Description)
	} else {
		description = types.StringNull()
	}

	model := &FirewallPolicyServerModel{
		Id:              types.StringValue(fp.Id),
		Name:            types.StringValue(fp.Name),
		Description:     description,
		State:           types.StringValue(fp.State),
		CreationDate:    types.StringValue(fp.CreationDate),
		Default:         types.Int64Value(int64(fp.Default)),
		Rules:           rulesList,
		ServerIPs:       serverIPsSet,
		ServerIPsDetail: serverIPsList,
		CloudPanelId:    types.StringValue(fp.CloudPanelID),
	}

	return model, diags
}

func (m *FirewallPolicyServerModel) ToAssignRequest(ctx context.Context) (*FirewallPolicyServerAssignRequest, diag.Diagnostics) {
	var serverIPs []string

	diags := m.ServerIPs.ElementsAs(ctx, &serverIPs, false)
	if diags.HasError() {
		return nil, diags
	}

	return &FirewallPolicyServerAssignRequest{
		ServerIPs: serverIPs,
	}, nil
}

type FirewallPolicyServerIpModel struct {
	FirewallPolicyId types.String `tfsdk:"firewall_policy_id"`
	ServerIpId       types.String `tfsdk:"server_ip_id"`
	IP               types.String `tfsdk:"ip"`
	ServerName       types.String `tfsdk:"server_name"`
}

func NewFirewallPolicyServerIpDataSourceModel(_ context.Context, firewallPolicyId string, response firewallPolicies.FirewallServerIPResponse) *FirewallPolicyServerIpModel {
	return &FirewallPolicyServerIpModel{
		FirewallPolicyId: types.StringValue(firewallPolicyId),
		ServerIpId:       types.StringValue(response.Id),
		IP:               types.StringValue(response.IP),
		ServerName:       types.StringValue(response.ServerName),
	}
}

func FirewallPolicyServerIpDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for firewall policy server Ip information",
		Attributes: map[string]schema.Attribute{
			"firewall_policy_id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"server_ip_id": schema.StringAttribute{
				Required:    true,
				Description: "Public IP identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "Server IP address",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.IPv4Pattern),
						"must be a valid IPv4 address",
					),
				},
			},
			"server_name": schema.StringAttribute{
				Computed:    true,
				Description: "Server name",
			},
		},
	}
}
func FirewallPolicyAssignmentResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Assigns server IPs to an existing firewall policy. Changes to firewall_policy_id or server_ips will force resource replacement (destroy + create).",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Description: "Firewall policy ID to assign server IPs to",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid Id (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"server_ips": rschema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Description: "List of server IP Ids to assign to the firewall policy",
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
				Description: "Firewall policy rules",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: rschema.NestedAttributeObject{
					Attributes: firewallPolicies.FirewallRuleResourceSchema(),
				},
			},
			"server_ips_detail": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "ServerIPs assigned to firewall policy",
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
