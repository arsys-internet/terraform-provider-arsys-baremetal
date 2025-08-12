package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type DatacentersListModel struct {
	ID          types.String `tfsdk:"id"`
	Datacenters types.List   `tfsdk:"datacenters"`
}

type DatacentersResponse = DatacenterResponse

func DatacentersDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing datacenters",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"datacenters": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of datacenters",
				NestedObject: datacenterNestedAttributeObject(),
			},
		},
	}
}

func NewDatacenters(ctx context.Context, datacentersResponse []DatacentersResponse) (*DatacentersListModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &DatacentersListModel{}
	model.ID = types.StringValue("datacenters")

	datacenterModels, listDiags := NewDatacenterFromList(ctx, datacentersResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		datacentersList, convertDiags := types.ListValueFrom(ctx, datacenterObjectType(), datacenterModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.Datacenters = datacentersList
		}
	}

	return model, diags
}
