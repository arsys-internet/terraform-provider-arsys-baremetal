package models

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalModels struct {
	Id              types.String `tfsdk:"id"`
	BaremetalModels types.List   `tfsdk:"baremetal_models"`
}

type BaremetalModelsResponse = BaremetalModelResponse

func NewBaremetalModelFromList(ctx context.Context, bmList []BaremetalModelResponse) ([]BaremetalModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []BaremetalModel

	if len(bmList) == 0 {
		return []BaremetalModel{}, diags
	}

	for i, bm := range bmList {
		model, modelDiags := NewBaremetalModel(ctx, bm)
		if modelDiags.HasError() {
			diags.AddError(
				"Build error",
				fmt.Sprintf("Failed to create model for item %d: %s", i, modelDiags.Errors()[0].Summary()),
			)
			continue
		}
		diags.Append(modelDiags...)

		models = append(models, model)

	}

	return models, diags
}

func NewBaremetalModels(ctx context.Context, datacentersResponse []BaremetalModelResponse) (*BaremetalModels, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &BaremetalModels{}
	model.Id = types.StringValue("baremetal_models")

	baremetalModels, listDiags := NewBaremetalModelFromList(ctx, datacentersResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		baremetalModelsList, convertDiags := types.ListValueFrom(ctx, BaremetalModelObjectType(), baremetalModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.BaremetalModels = baremetalModelsList
		}
	}

	return model, diags
}

func baremetalModelNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := BaremetalModelDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Baremetal model identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func BaremetalModelsDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing baremetal models",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"baremetal_models": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of baremetal models",
				NestedObject: baremetalModelNestedAttributeObject(),
			},
		},
	}
}
