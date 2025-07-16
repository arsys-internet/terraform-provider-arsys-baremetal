package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type StatusModel struct {
	State   types.String `tfsdk:"state"`
	Percent types.Int64  `tfsdk:"percent"`
}

func NewStatusModel(status StatusResponse) StatusModel {
	return StatusModel{
		State:   types.StringValue(status.State),
		Percent: types.Int64Value(int64(status.Percent)),
	}
}

type StatusResponse struct {
	State   string `json:"state"`
	Percent int    `json:"percent"`
}

func NewStatusObject(status StatusResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"state":   types.StringValue(status.State),
		"percent": types.Int64Value(int64(status.Percent)),
	}
	return types.ObjectValue(StatusObjectType().AttrTypes, attrs)
}

func StatusObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"state":   types.StringType,
			"percent": types.Int64Type,
		},
	}
}

func StatusDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"state": schema.StringAttribute{
			Computed:    true,
			Description: "Current state of the server",
		},
		"percent": schema.Int64Attribute{
			Computed:    true,
			Description: "Percentage of completion for the current operation",
		},
	}
}

func StatusResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"state": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Current state of the server",
		},
		"percent": rschema.Int64Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Description: "Percentage of completion for the current operation",
		},
	}
}
