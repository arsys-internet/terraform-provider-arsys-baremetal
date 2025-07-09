package models

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ServersModel struct {
	ID      types.String `tfsdk:"id"`
	Servers types.List   `tfsdk:"servers"`
}

func NewServers(ctx context.Context, serversResponse []ServerResponse) (*ServersModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &ServersModel{}
	model.ID = types.StringValue("servers")

	serverModels, listDiags := NewServerFromList(ctx, serversResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		serversList, convertDiags := types.ListValueFrom(ctx, serverModelObjectType(), serverModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.Servers = serversList
		}
	}

	return model, diags
}

func ServersDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing servers",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"servers": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers",
				NestedObject: serverNestedAttributeObject(),
			},
		},
	}
}
