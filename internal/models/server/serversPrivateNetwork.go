package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServersPrivateNetworkResponse struct {
	Id       string  `json:"id"`
	Name     string  `json:"name"`
	ServerIP string  `json:"server_ip"`
	VlanId   *string `json:"vlan_id"`
}

func NewServersPrivateNetworkObject(pn ServersPrivateNetworkResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"id":        types.StringValue(pn.Id),
		"name":      types.StringValue(pn.Name),
		"server_ip": types.StringValue(pn.ServerIP),
	}

	if pn.VlanId != nil {
		attrs["vlan_id"] = types.StringValue(*pn.VlanId)
	} else {
		attrs["vlan_id"] = types.StringNull()
	}

	return types.ObjectValue(ServersPrivateNetworkObjectType().AttrTypes, attrs)
}

func NewServersPrivateNetworkList(privateNetworks []ServersPrivateNetworkResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(privateNetworks) == 0 {
		return types.ListValueMust(ServersPrivateNetworkObjectType(), []attr.Value{}), diags
	}

	elements := make([]attr.Value, 0, len(privateNetworks))

	for _, pn := range privateNetworks {
		pnObj, objDiags := NewServersPrivateNetworkObject(pn)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			elements = append(elements, pnObj)
		}
	}

	if diags.HasError() {
		return types.ListValueMust(ServersPrivateNetworkObjectType(), []attr.Value{}), diags
	}

	list, listDiags := types.ListValue(ServersPrivateNetworkObjectType(), elements)
	diags.Append(listDiags...)

	return list, diags
}

func ServersPrivateNetworkObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":        types.StringType,
			"name":      types.StringType,
			"server_ip": types.StringType,
			"vlan_id":   types.StringType,
		},
	}
}

func ServersPrivateNetworksDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Private network identifier",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Private network name",
		},
		"server_ip": schema.StringAttribute{
			Computed:    true,
			Description: "Server IP address in the private network",
		},
		"vlan_id": schema.StringAttribute{
			Computed:    true,
			Description: "VLAN identifier for baremetal servers",
		},
	}
}

func ServersPrivateNetworksResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Private network identifier",
		},
		"name": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Private network name",
		},
		"server_ip": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Server IP address in the private network",
		},
		"vlan_id": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "VLAN identifier for baremetal servers",
		},
	}
}
