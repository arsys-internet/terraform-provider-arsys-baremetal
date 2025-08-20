package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AvailabilityResponse struct {
	DatacenterId               string `json:"datacenter_id"`
	Available                  bool   `json:"available"`
	AvailableConnectionsSpeeds []int  `json:"available_connections_speeds"`
	AvailableWithRedundancy    bool   `json:"available_with_redundancy"`
	AvailableWithoutRedundancy bool   `json:"available_without_redundancy"`
}

type AvailabilityModel struct {
	DatacenterId               types.String `tfsdk:"datacenter_id"`
	Available                  types.Bool   `tfsdk:"available"`
	AvailableConnectionsSpeeds types.List   `tfsdk:"available_connections_speeds"`
	AvailableWithRedundancy    types.Bool   `tfsdk:"available_with_redundancy"`
	AvailableWithoutRedundancy types.Bool   `tfsdk:"available_without_redundancy"`
}

func NewAvailabilityObject(availability AvailabilityResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]attr.Value, 0, len(availability.AvailableConnectionsSpeeds))
	for _, speed := range availability.AvailableConnectionsSpeeds {
		elements = append(elements, types.Int64Value(int64(speed)))
	}

	speedsList, listDiags := types.ListValue(types.Int64Type, elements)
	diags.Append(listDiags...)

	availabilityObj, objDiags := types.ObjectValue(AvailabilityObjectType().AttrTypes,
		map[string]attr.Value{
			"datacenter_id":                types.StringValue(availability.DatacenterId),
			"available":                    types.BoolValue(availability.Available),
			"available_connections_speeds": speedsList,
			"available_with_redundancy":    types.BoolValue(availability.AvailableWithRedundancy),
			"available_without_redundancy": types.BoolValue(availability.AvailableWithoutRedundancy),
		})
	diags.Append(objDiags...)

	return availabilityObj, diags
}

func AvailabilityObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"datacenter_id":                types.StringType,
			"available":                    types.BoolType,
			"available_connections_speeds": types.ListType{ElemType: types.Int64Type},
			"available_with_redundancy":    types.BoolType,
			"available_without_redundancy": types.BoolType,
		},
	}
}

func AvailabilityDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"datacenter_id": schema.StringAttribute{
			Computed:    true,
			Description: "Datacenter identifier",
		},
		"available": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether the server is available",
		},
		"available_connections_speeds": schema.ListAttribute{
			ElementType: types.Int64Type,
			Computed:    true,
			Description: "List of available connection speeds",
		},
		"available_with_redundancy": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether available with redundancy",
		},
		"available_without_redundancy": schema.BoolAttribute{
			Computed:    true,
			Description: "Whether available without redundancy",
		},
	}
}
