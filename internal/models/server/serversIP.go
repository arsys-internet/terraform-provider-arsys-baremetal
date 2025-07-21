package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServersIPResponse struct {
	ID             string                 `json:"id"`
	IP             string                 `json:"ip"`
	Type           string                 `json:"type"`
	ReverseDNS     interface{}            `json:"reverse_dns,omitempty"`
	Main           bool                   `json:"main"`
	FirewallPolicy IdentifierIPResponse   `json:"firewall_policy"`
	LoadBalancers  []IdentifierIPResponse `json:"load_balancers"`
}

type IdentifierIPResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func NewServersIPList(ips []ServersIPResponse) (types.List, diag.Diagnostics) {
	var elements []attr.Value

	for _, ip := range ips {
		fpAttrs := map[string]attr.Value{
			"id":   types.StringValue(ip.FirewallPolicy.ID),
			"name": types.StringValue(ip.FirewallPolicy.Name),
		}
		fpObj, _ := types.ObjectValue(IdentifierIPObjectType().AttrTypes, fpAttrs)

		var lbElements []attr.Value
		for _, lb := range ip.LoadBalancers {
			lbAttrs := map[string]attr.Value{
				"id":   types.StringValue(lb.ID),
				"name": types.StringValue(lb.Name),
			}
			lbObj, _ := types.ObjectValue(IdentifierIPObjectType().AttrTypes, lbAttrs)
			lbElements = append(lbElements, lbObj)
		}
		lbList, _ := types.ListValue(IdentifierIPObjectType(), lbElements)

		reverseDNS := types.StringNull()
		if ip.ReverseDNS != nil {
			if str, ok := ip.ReverseDNS.(string); ok {
				reverseDNS = types.StringValue(str)
			}
		}

		ipAttrs := map[string]attr.Value{
			"id":              types.StringValue(ip.ID),
			"ip":              types.StringValue(ip.IP),
			"type":            types.StringValue(ip.Type),
			"reverse_dns":     reverseDNS,
			"main":            types.BoolValue(ip.Main),
			"firewall_policy": fpObj,
			"load_balancers":  lbList,
		}

		ipObj, _ := types.ObjectValue(ServersIPObjectType().AttrTypes, ipAttrs)
		elements = append(elements, ipObj)
	}

	return types.ListValue(ServersIPObjectType(), elements)
}

func ServersIPObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":              types.StringType,
			"ip":              types.StringType,
			"type":            types.StringType,
			"reverse_dns":     types.StringType,
			"main":            types.BoolType,
			"firewall_policy": IdentifierIPObjectType(),
			"load_balancers":  types.ListType{ElemType: IdentifierIPObjectType()},
		},
	}
}

func IdentifierIPObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
		},
	}
}

func ServersIPDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "IP identifier",
		},
		"ip": schema.StringAttribute{
			Computed:    true,
			Description: "IP",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: "IP type",
		},
		"reverse_dns": schema.StringAttribute{
			Computed:    true,
			Description: "Reverse name of the IP",
		},
		"main": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether this is the main IP",
		},
		"firewall_policy": schema.SingleNestedAttribute{
			Computed:    true,
			Description: "Firewall policy assigned to IP",
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Computed:    true,
					Description: "Firewall policy ID",
				},
				"name": schema.StringAttribute{
					Computed:    true,
					Description: "Firewall policy name",
				},
			},
		},
		"load_balancers": schema.ListNestedAttribute{
			Computed:    true,
			Description: "Load balancer(s) assigned to IP",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "Load balancer ID",
					},
					"name": schema.StringAttribute{
						Computed:    true,
						Description: "Load balancer name",
					},
				},
			},
		},
	}
}

func ServersIPResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "IP identifier",
		},
		"ip": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "IP",
		},
		"type": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "IP type",
		},
		"reverse_dns": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Reverse name of the IP",
		},
		"main": rschema.BoolAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
			Description: "Whether this is the main IP",
		},
		"firewall_policy": rschema.SingleNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Description: "Firewall policy assigned to IP",
			Attributes: map[string]rschema.Attribute{
				"id": rschema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Description: "Firewall policy ID",
				},
				"name": rschema.StringAttribute{
					Computed: true,
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
					Description: "Firewall policy name",
				},
			},
		},
		"load_balancers": rschema.ListNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description: "Load balancer(s) assigned to IP",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"id": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Load balancer ID",
					},
					"name": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Load balancer name",
					},
				},
			},
		},
	}
}
