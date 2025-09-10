package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetworkServersModel struct {
	Id      types.String `tfsdk:"id"`
	Servers types.List   `tfsdk:"servers"`
}

func NewPrivateNetworkServers(ctx context.Context, id string, privateNetworkServersResponse []PrivateNetworkServerResponse) (*PrivateNetworkServersModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &PrivateNetworkServersModel{}
	model.Id = types.StringValue(id)
	privateNetworkServers, listDiags := NewPrivateNetworkServersList(ctx, privateNetworkServersResponse)
	diags.Append(diags...)

	if !listDiags.HasError() {
		privateNetworkServersList, convertDiags := types.ListValueFrom(ctx, privateNetworkServerObjectType(), privateNetworkServers)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.Servers = privateNetworkServersList
		}
	}

	return model, diags
}

func NewPrivateNetworkServerListObject(server PrivateNetworkServerResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var lockValue types.Int64
	if server.Lock != nil {
		lockValue = types.Int64Value(*server.Lock)
	} else {
		lockValue = types.Int64Null()
	}
	attrs := map[string]attr.Value{
		"id":   types.StringValue(server.Id),
		"name": types.StringValue(server.Name),
		"lock": lockValue,
	}

	obj, objDiags := types.ObjectValue(privateNetworkServerObjectType().AttrTypes, attrs)
	diags.Append(objDiags...)

	return obj, diags
}

func NewPrivateNetworkServersList(_ context.Context, servers []PrivateNetworkServerResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(servers) == 0 {
		return types.ListValueMust(privateNetworkServerObjectType(), []attr.Value{}), diags
	}

	elements := make([]attr.Value, 0, len(servers))

	for _, server := range servers {
		serverObj, objDiags := NewPrivateNetworkServerListObject(server)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			elements = append(elements, serverObj)
		}
	}

	if diags.HasError() {
		return types.ListValueMust(privateNetworkServerObjectType(), []attr.Value{}), diags
	}

	list, listDiags := types.ListValue(privateNetworkServerObjectType(), elements)
	diags.Append(listDiags...)

	return list, diags
}

func PrivateNetworkServersDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing private network servers",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Private network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid private network ID",
					),
				},
			},
			"servers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of private network servers",
				NestedObject: schema.NestedAttributeObject{
					Attributes: PrivateNetworkServerListDataSourceSchema(),
				},
			},
		},
	}
}

func PrivateNetworkServerListDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Private network server identifier",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Private network server name",
		},
		"lock": schema.Int64Attribute{
			Computed:    true,
			Description: "Private network server lock status",
		},
	}
}
