package models

import (
	"context"
	"terraform-provider-arsys-baremetal/internal/models/server"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalModelResponse struct {
	Id           string                           `json:"id"`
	Name         string                           `json:"name"`
	Hardware     server.BaremetalHardwareResponse `json:"hardware"`
	StateId      int                              `json:"state_id"`
	State        string                           `json:"state"`
	Availability []server.AvailabilityResponse    `json:"availability"`
}

type BaremetalModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Hardware     types.Object `tfsdk:"hardware"`
	StateId      types.Int64  `tfsdk:"state_id"`
	State        types.String `tfsdk:"state"`
	Availability types.List   `tfsdk:"availability"`
}

func NewBaremetalModel(_ context.Context, bmResponse BaremetalModelResponse) (BaremetalModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := BaremetalModel{
		Id:      types.StringValue(bmResponse.Id),
		Name:    types.StringValue(bmResponse.Name),
		StateId: types.Int64Value(int64(bmResponse.StateId)),
		State:   types.StringValue(bmResponse.State),
	}

	hardwareObj, hardwareDiags := server.NewBaremetalHardwareObject(bmResponse.Hardware)
	diags.Append(hardwareDiags...)
	if !hardwareDiags.HasError() {
		model.Hardware = hardwareObj
	}

	availabilityElements := make([]attr.Value, 0, len(bmResponse.Availability))
	for _, availability := range bmResponse.Availability {
		availabilityObj, availabilityDiags := server.NewAvailabilityObject(availability)
		diags.Append(availabilityDiags...)

		if !availabilityDiags.HasError() {
			availabilityElements = append(availabilityElements, availabilityObj)
		}
	}

	availabilityList, listDiags := types.ListValue(server.AvailabilityObjectType(), availabilityElements)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.Availability = availabilityList
	}

	return model, diags
}

func BaremetalModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":           types.StringType,
			"name":         types.StringType,
			"hardware":     server.BaremetalHardwareObjectType(),
			"state_id":     types.Int64Type,
			"state":        types.StringType,
			"availability": types.ListType{ElemType: server.AvailabilityObjectType()},
		},
	}
}

func BaremetalModelDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Fetches information about a specific baremetal model",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Baremetal model identifier",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Baremetal model name",
			},
			"hardware": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Hardware specifications",
				Attributes:  server.BaremetalHardwareDataSourceSchema(),
			},
			"state_id": schema.Int64Attribute{
				Computed:    true,
				Description: "State identifier",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "State description",
			},
			"availability": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of availability per datacenter",
				NestedObject: schema.NestedAttributeObject{
					Attributes: server.AvailabilityDataSourceSchema(),
				},
			},
		},
	}
}
