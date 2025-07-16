package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type RedundancyResponse struct {
	Available bool `json:"available"`
	Enabled   bool `json:"enabled"`
}

func NewRedundancyObject(redundancy RedundancyResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"available": types.BoolValue(redundancy.Available),
		"enabled":   types.BoolValue(redundancy.Enabled),
	}
	return types.ObjectValue(RedundancyObjectType().AttrTypes, attrs)
}

func RedundancyObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"available": types.BoolType,
			"enabled":   types.BoolType,
		},
	}
}

func RedundancyDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"available": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether redundancy is available for this server",
		},
		"enabled": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether redundancy is currently enabled",
		},
	}
}

func RedundancyResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"available": rschema.BoolAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
			Description: "Whether redundancy is available for this server",
		},
		"enabled": rschema.BoolAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.Bool{
				boolplanmodifier.UseStateForUnknown(),
			},
			Description: "Whether redundancy is currently enabled",
		},
	}
}
