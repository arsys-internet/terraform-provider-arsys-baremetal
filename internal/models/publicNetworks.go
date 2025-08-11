package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicNetworksModel struct {
	Id             types.String `tfsdk:"id"`
	PublicNetworks types.List   `tfsdk:"public_networks"`
}

func NewPublicNetworks(ctx context.Context, publicNetworksResponse []PublicNetworkResponse) (*PublicNetworksModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &PublicNetworksModel{}
	model.Id = types.StringValue("public_networks")

	publicNetworkModels, listDiags := NewPublicNetworkFromList(ctx, publicNetworksResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		publicNetworksList, convertDiags := types.ListValueFrom(ctx, publicNetworkObjectType(), publicNetworkModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.PublicNetworks = publicNetworksList
		}
	}

	return model, diags
}

func PublicNetworksDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing public networks",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"public_networks": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of public networks",
				NestedObject: publicNetworkNestedAttributeObject(),
			},
		},
	}
}
