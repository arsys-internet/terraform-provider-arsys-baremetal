package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type AlertResponse struct {
	Critical []AlertItemResponse `json:"critical"`
	Warning  []AlertItemResponse `json:"warning"`
}

type AlertItemResponse struct {
	Date        string `json:"date"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func NewAlertsObject(alerts *AlertResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if alerts == nil {
		return types.ObjectNull(AlertsObjectType().AttrTypes), diags
	}

	criticalElements := make([]attr.Value, 0, len(alerts.Critical))
	for _, alert := range alerts.Critical {
		alertObj, objDiags := NewAlertItemObject(alert)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			criticalElements = append(criticalElements, alertObj)
		}
	}

	warningElements := make([]attr.Value, 0, len(alerts.Warning))
	for _, alert := range alerts.Warning {
		alertObj, objDiags := NewAlertItemObject(alert)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			warningElements = append(warningElements, alertObj)
		}
	}

	criticalList, criticalDiags := types.ListValue(AlertItemObjectType(), criticalElements)
	diags.Append(criticalDiags...)

	warningList, warningDiags := types.ListValue(AlertItemObjectType(), warningElements)
	diags.Append(warningDiags...)

	alertsObj, objDiags := types.ObjectValue(AlertsObjectType().AttrTypes, map[string]attr.Value{
		"critical": criticalList,
		"warning":  warningList,
	})
	diags.Append(objDiags...)

	return alertsObj, diags
}

func NewAlertItemObject(alert AlertItemResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"date":        types.StringValue(alert.Date),
		"description": types.StringValue(alert.Description),
		"type":        types.StringValue(alert.Type),
	}

	return types.ObjectValue(AlertItemObjectType().AttrTypes, attrs)
}

func AlertsObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"critical": types.ListType{ElemType: AlertItemObjectType()},
			"warning":  types.ListType{ElemType: AlertItemObjectType()},
		},
	}
}

func AlertItemObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"date":        types.StringType,
			"description": types.StringType,
			"type":        types.StringType,
		},
	}
}

func AlertsDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"critical": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of critical alerts",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"date": schema.StringAttribute{
						Computed:    true,
						Description: "Alert date",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Alert description",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Alert type",
					},
				},
			},
		},
		"warning": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of warning alerts",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"date": schema.StringAttribute{
						Computed:    true,
						Description: "Alert date",
					},
					"description": schema.StringAttribute{
						Computed:    true,
						Description: "Alert description",
					},
					"type": schema.StringAttribute{
						Computed:    true,
						Description: "Alert type",
					},
				},
			},
		},
	}
}

func AlertsResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"critical": rschema.ListNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description: "List of critical alerts",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"date": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Alert date",
					},
					"description": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Alert description",
					},
					"type": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Alert type",
					},
				},
			},
		},
		"warning": rschema.ListNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description: "List of warning alerts",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"date": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Alert date",
					},
					"description": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Alert description",
					},
					"type": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Alert type",
					},
				},
			},
		},
	}
}
