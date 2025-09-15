package models

import (
	"context"
	"fmt"
	"regexp"
	firewallpolicy "terraform-provider-arsys-baremetal/internal/models/firewall_policy"

	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPolicyServerIpsModel struct {
	Id        types.String `tfsdk:"id"`
	ServerIPs types.List   `tfsdk:"server_ips"`
}

func NewFirewallPolicyServerIpsModel(_ context.Context, id string, servers []firewallpolicy.FirewallServerIPResponse) (*FirewallPolicyServerIpsModel, error) {
	serverIPsList, diags := firewallpolicy.NewFirewallServerIPsList(servers)
	if diags.HasError() {
		return nil, fmt.Errorf("failed to create server IPs list: %v", diags)
	}

	return &FirewallPolicyServerIpsModel{
		Id:        types.StringValue(id),
		ServerIPs: serverIPsList,
	}, nil
}

func FirewallPolicyServerIPsSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Firewall policy servers",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Id of the firewall policy",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid Id",
					),
				},
			},
			"server_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "ServerIPs assigned to firewall policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: firewallpolicy.FirewallServerIPDataSourceSchema(),
				},
			},
		},
	}
}
