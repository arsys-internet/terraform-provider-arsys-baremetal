package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServerAppliancesModel struct {
	Id               types.String `tfsdk:"id"`
	ServerAppliances types.List   `tfsdk:"server_appliances"`
}

func NewServerAppliances(ctx context.Context, serverAppliancesResponse []ServerApplianceResponse) (*ServerAppliancesModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &ServerAppliancesModel{}
	model.Id = types.StringValue("server_appliances")

	serverApplianceModels, listDiags := NewServerApplianceFromList(ctx, serverAppliancesResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		serverAppliancesList, convertDiags := types.ListValueFrom(ctx, serverApplianceObjectType(), serverApplianceModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.ServerAppliances = serverAppliancesList
		}
	}

	return model, diags
}

func ServerAppliancesDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing server appliances",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"server_appliances": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of server appliances",
				NestedObject: serverApplianceNestedAttributeObject(),
			},
		},
	}
}
