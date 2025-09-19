package firewallpolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallServerIPResponse struct {
	Id         string `json:"id"`
	IP         string `json:"ip"`
	ServerName string `json:"server_name"`
}

func NewFirewallServerIPObject(serverIP FirewallServerIPResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"id":          types.StringValue(serverIP.Id),
		"ip":          types.StringValue(serverIP.IP),
		"server_name": types.StringValue(serverIP.ServerName),
	}

	return types.ObjectValue(FirewallServerIPObjectType().AttrTypes, attrs)
}

func NewFirewallServerIPsList(serverIPs []FirewallServerIPResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(serverIPs) == 0 {
		return types.ListValueMust(FirewallServerIPObjectType(), []attr.Value{}), diags
	}

	elements := make([]attr.Value, 0, len(serverIPs))

	for _, serverIP := range serverIPs {
		serverIPObj, objDiags := NewFirewallServerIPObject(serverIP)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			elements = append(elements, serverIPObj)
		}
	}

	if diags.HasError() {
		return types.ListValueMust(FirewallServerIPObjectType(), []attr.Value{}), diags
	}

	list, listDiags := types.ListValue(FirewallServerIPObjectType(), elements)
	diags.Append(listDiags...)

	return list, diags
}

func FirewallServerIPObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":          types.StringType,
			"ip":          types.StringType,
			"server_name": types.StringType,
		},
	}
}

func FirewallServerIPDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Server identifier",
		},
		"ip": schema.StringAttribute{
			Computed:    true,
			Description: "Server IP address",
		},
		"server_name": schema.StringAttribute{
			Computed:    true,
			Description: "Server name",
		},
	}
}

func FirewallServerIPResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed:    true,
			Description: "Server identifier",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"ip": rschema.StringAttribute{
			Computed:    true,
			Description: "Server IP address",
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
	}
}
