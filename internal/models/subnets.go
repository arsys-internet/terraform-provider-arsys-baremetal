package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type SubnetsModel struct {
	Id      types.String `tfsdk:"id"`
	Subnets types.List   `tfsdk:"subnets"`
}

func NewSubnets(ctx context.Context, publicIpsResponse []PublicIpResponse) (*SubnetsModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &SubnetsModel{}
	model.Id = types.StringValue("subnets")

	publicIpModels, listDiags := NewPublicIpFromList(ctx, publicIpsResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		publicIpsList, convertDiags := types.ListValueFrom(ctx, publicIpObjectType(), publicIpModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.Subnets = publicIpsList
		}
	}

	return model, diags
}

func SubnetsDatasourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing subnets",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"subnets": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of subnets",
				NestedObject: publicIpNestedAttributeObject(),
			},
		},
	}
}
