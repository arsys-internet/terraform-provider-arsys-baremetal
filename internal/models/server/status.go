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

type StatusDetailModel struct {
	State   types.String `tfsdk:"state"`
	Percent types.Int64  `tfsdk:"percent"`
}

type StatusBaseModel struct {
	State types.String `tfsdk:"state"`
}

type StatusDetailResponse struct {
	State   string `json:"state"`
	Percent *int   `json:"percent"`
}
type StatusBaseResponse struct {
	State string `json:"state"`
}

func NewStatusDetailObject(status StatusDetailResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"state": types.StringValue(status.State),
	}

	if status.Percent != nil {
		attrs["percent"] = types.Int64Value(int64(*status.Percent))
	} else {
		attrs["percent"] = types.Int64Null()
	}

	return types.ObjectValue(StatusDetailObjectType().AttrTypes, attrs)
}

func NewStatusBaseObject(status StatusBaseResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"state": types.StringValue(status.State),
	}
	return types.ObjectValue(StatusBaseObjectType().AttrTypes, attrs)
}

func StatusDetailObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"state":   types.StringType,
			"percent": types.Int64Type,
		},
	}
}

func StatusBaseObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"state": types.StringType,
		},
	}
}

func StatusDetailDataSourceSchema() map[string]schema.Attribute {
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

func StatusDetailResourceSchema() map[string]rschema.Attribute {
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
