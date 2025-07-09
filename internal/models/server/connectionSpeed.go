package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ConnectionSpeedResponse struct {
	// For servers cloud (estructura simple)
	Available []float64 `json:"available,omitempty"`
	Current   *float64  `json:"current,omitempty"`

	// For servers baremetal (estructura anidada)
	Private *ConnectionSpeedDetailResponse `json:"private,omitempty"`
	Public  *ConnectionSpeedDetailResponse `json:"public,omitempty"`
}

type ConnectionSpeedDetailResponse struct {
	Available []float64 `json:"available"`
	Current   float64   `json:"current"`
}

func NewConnectionSpeedObject(cs ConnectionSpeedResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{}

	// For cloud servers (estructura simple)
	if len(cs.Available) > 0 {
		availableElements := make([]attr.Value, len(cs.Available))
		for i, v := range cs.Available {
			availableElements[i] = types.Float64Value(v)
		}
		availableList, _ := types.ListValue(types.Float64Type, availableElements)
		attrs["available"] = availableList
	} else {
		attrs["available"] = types.ListNull(types.Float64Type)
	}

	if cs.Current != nil {
		attrs["current"] = types.Float64Value(*cs.Current)
	} else {
		attrs["current"] = types.Float64Null()
	}

	// For baremetal servers (estructura anidada)
	if cs.Private != nil {
		privateAttrs := map[string]attr.Value{}

		privAvailableElements := make([]attr.Value, len(cs.Private.Available))
		for i, v := range cs.Private.Available {
			privAvailableElements[i] = types.Float64Value(v)
		}
		privAvailableList, _ := types.ListValue(types.Float64Type, privAvailableElements)
		privateAttrs["available"] = privAvailableList
		privateAttrs["current"] = types.Float64Value(cs.Private.Current)

		privateObj, _ := types.ObjectValue(ConnectionSpeedDetailObjectType().AttrTypes, privateAttrs)
		attrs["private"] = privateObj
	} else {
		attrs["private"] = types.ObjectNull(ConnectionSpeedDetailObjectType().AttrTypes)
	}

	if cs.Public != nil {
		publicAttrs := map[string]attr.Value{}

		pubAvailableElements := make([]attr.Value, len(cs.Public.Available))
		for i, v := range cs.Public.Available {
			pubAvailableElements[i] = types.Float64Value(v)
		}
		pubAvailableList, _ := types.ListValue(types.Float64Type, pubAvailableElements)
		publicAttrs["available"] = pubAvailableList
		publicAttrs["current"] = types.Float64Value(cs.Public.Current)

		publicObj, _ := types.ObjectValue(ConnectionSpeedDetailObjectType().AttrTypes, publicAttrs)
		attrs["public"] = publicObj
	} else {
		attrs["public"] = types.ObjectNull(ConnectionSpeedDetailObjectType().AttrTypes)
	}

	return types.ObjectValue(ConnectionSpeedObjectType().AttrTypes, attrs)
}

func ConnectionSpeedObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			// For cloud servers
			"available": types.ListType{ElemType: types.Float64Type},
			"current":   types.Float64Type,
			// For baremetal servers
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
		"available": schema.ListAttribute{
			ElementType: types.Float64Type,
			Computed:    true,
			Description: "Available connection speeds for cloud servers",
		},
		"current": schema.Float64Attribute{
			Computed:    true,
			Description: "Current connection speed for cloud servers",
		},
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
		"available": rschema.ListAttribute{
			ElementType: types.Float64Type,
			Computed:    true,
			Description: "Available connection speeds for cloud servers",
		},
		"current": rschema.Float64Attribute{
			Computed:    true,
			Description: "Current connection speed for cloud servers",
		},
		"private": rschema.SingleNestedAttribute{
			Attributes: map[string]rschema.Attribute{
				"available": rschema.ListAttribute{
					ElementType: types.Float64Type,
					Computed:    true,
					Description: "Available private connection speeds for baremetal servers",
				},
				"current": rschema.Float64Attribute{
					Computed:    true,
					Description: "Current private connection speed for baremetal servers",
				},
			},
			Computed:    true,
			Description: "Private connection speed details for baremetal servers",
		},
		"public": rschema.SingleNestedAttribute{
			Attributes: map[string]rschema.Attribute{
				"available": rschema.ListAttribute{
					ElementType: types.Float64Type,
					Computed:    true,
					Description: "Available public connection speeds for baremetal servers",
				},
				"current": rschema.Float64Attribute{
					Computed:    true,
					Description: "Current public connection speed for baremetal servers",
				},
			},
			Computed:    true,
			Description: "Public connection speed details for baremetal servers",
		},
	}
}

// ... resto de tu código existente ...
