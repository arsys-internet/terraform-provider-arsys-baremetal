package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetworksModel struct {
	ID              types.String `tfsdk:"id"`
	PrivateNetworks types.List   `tfsdk:"private_networks"`
}

func NewPrivateNetworks(ctx context.Context, privateNetworksResponse []PrivateNetworkResponse) (*PrivateNetworksModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	model := &PrivateNetworksModel{}
	model.ID = types.StringValue("private_networks")

	privateNetworkModels, listDiags := NewPrivateNetworkFromList(ctx, privateNetworksResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		privateNetworksList, convertDiags := types.ListValueFrom(ctx, privateNetworkObjectType(), privateNetworkModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.PrivateNetworks = privateNetworksList
		}
	}

	return model, diags
}

func PrivateNetworksDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing private networks",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"private_networks": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of private networks",
				NestedObject: privateNetworkNestedAttributeObject(),
			},
		},
	}
}
