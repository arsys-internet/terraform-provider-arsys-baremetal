package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier" // AÑADIR ESTA LÍNEA
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/firewallPolicies"
	"terraform-provider-arsys-baremetal/internal/util"
)

type FirewallPolicyServerRequest struct {
	Servers []string `json:"servers"`
}

type FirewallPolicyServerDeleteRequest struct {
	ServerIpId string `json:"server_ip"`
}

type FirewallPolicyServerResourceModel struct {
	PublicNetworkId types.String   `tfsdk:"public_network_id"`
	Servers         []types.String `tfsdk:"servers"`
	Id              types.String   `tfsdk:"id"`
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
		ServerIpId:       types.StringValue(response.ID),
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
			"ip": rschema.StringAttribute{
				Computed:    true,
				Description: "Server IP address",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.IPv4Pattern),
						"must be a valid IPv4 address",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_name": rschema.StringAttribute{
				Computed:    true,
				Description: "Server name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func FirewallPolicyAssignServerIPSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Firewall policy server ids resource",
		Attributes: map[string]rschema.Attribute{
			"firewall_policy_id": rschema.StringAttribute{
				Required:    true,
				Description: "Firewall policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"servers": rschema.ListAttribute{
				Required:    true,
				Description: "List of servers identifiers in the firwall policy",
				ElementType: types.StringType,
			},
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Internal ID for the resource",
			},
		},
	}
}

func FirewallPolicyServerDeleteSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Firewall policy server ids resource",
		Attributes: map[string]rschema.Attribute{
			"firewall_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Firewall policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"server_ip_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Server Ip identifiers in the firewall policy",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Internal ID for the resource",
			},
		},
	}
}
