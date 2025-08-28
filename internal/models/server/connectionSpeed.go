package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConnectionSpeedResponse struct {
	Private ConnectionSpeedDetailResponse `json:"private"`
	Public  ConnectionSpeedDetailResponse `json:"public"`
}

type ConnectionSpeedDetailResponse struct {
	Available []float64 `json:"available"`
	Current   float64   `json:"current"`
}

func NewConnectionSpeedObject(cs ConnectionSpeedResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	attrs := map[string]attr.Value{}

	privateObj, privateDiags := newConnectionSpeedDetailObject(cs.Private)
	diags.Append(privateDiags...)
	attrs["private"] = privateObj

	publicObj, publicDiags := newConnectionSpeedDetailObject(cs.Public)
	diags.Append(publicDiags...)
	attrs["public"] = publicObj

	obj, objDiags := types.ObjectValue(ConnectionSpeedObjectType().AttrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}

func newConnectionSpeedDetailObject(detail ConnectionSpeedDetailResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	availableElements := make([]attr.Value, len(detail.Available))
	for i, v := range detail.Available {
		availableElements[i] = types.Float64Value(v)
	}

	availableList, listDiags := types.ListValue(types.Float64Type, availableElements)
	diags.Append(listDiags...)

	attrs := map[string]attr.Value{
		"available": availableList,
		"current":   types.Float64Value(detail.Current),
	}

	obj, objDiags := types.ObjectValue(ConnectionSpeedDetailObjectType().AttrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}

func ConnectionSpeedObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"private": ConnectionSpeedDetailObjectType(),
			"public":  ConnectionSpeedDetailObjectType(),
		},
	}
}

func ConnectionSpeedDetailObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"available": types.ListType{ElemType: types.Float64Type},
			"current":   types.Float64Type,
		},
	}
}

func ConnectionSpeedDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"private": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"available": schema.ListAttribute{
					ElementType: types.Float64Type,
					Computed:    true,
					Description: "Available private connection speeds for baremetal servers",
				},
				"current": schema.Float64Attribute{
					Computed:    true,
					Description: "Current private connection speed for baremetal servers",
				},
			},
			Computed:    true,
			Description: "Private connection speed details for baremetal servers",
		},
		"public": schema.SingleNestedAttribute{
			Attributes: map[string]schema.Attribute{
				"available": schema.ListAttribute{
					ElementType: types.Float64Type,
					Computed:    true,
					Description: "Available public connection speeds for baremetal servers",
				},
				"current": schema.Float64Attribute{
					Computed:    true,
					Description: "Current public connection speed for baremetal servers",
				},
			},
			Computed:    true,
			Description: "Public connection speed details for baremetal servers",
		},
	}
}

func ConnectionSpeedResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"private": rschema.SingleNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]rschema.Attribute{
				"available": rschema.ListAttribute{
					ElementType: types.Float64Type,
					Computed:    true,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
					Description: "Available private connection speeds for baremetal servers",
				},
				"current": rschema.Float64Attribute{
					Computed: true,
					PlanModifiers: []planmodifier.Float64{
						float64planmodifier.UseStateForUnknown(),
					},
					Description: "Current private connection speed for baremetal servers",
				},
			},
			Description: "Private connection speed details for baremetal servers",
		},
		"public": rschema.SingleNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Object{
				objectplanmodifier.UseStateForUnknown(),
			},
			Attributes: map[string]rschema.Attribute{
				"available": rschema.ListAttribute{
					ElementType: types.Float64Type,
					Computed:    true,
					PlanModifiers: []planmodifier.List{
						listplanmodifier.UseStateForUnknown(),
					},
					Description: "Available public connection speeds for baremetal servers",
				},
				"current": rschema.Float64Attribute{
					Computed: true,
					PlanModifiers: []planmodifier.Float64{
						float64planmodifier.UseStateForUnknown(),
					},
					Description: "Current public connection speed for baremetal servers",
				},
			},
			Description: "Public connection speed details for baremetal servers",
		},
	}
}
