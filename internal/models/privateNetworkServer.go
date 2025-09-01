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

type PrivateNetworkServerModel struct {
	PrivateNetworkId types.String `tfsdk:"private_network_id"`
	Id               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Lock             types.Int64  `tfsdk:"lock"`
}

type PrivateNetworkServerResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Lock *int64 `json:"lock"`
}

type PrivateNetworkServerRequest struct {
	Servers []string `json:"servers"`
}

func NewPrivateNetworkServerFromResponse(_ context.Context, privateNetworkId string, response *PrivateNetworkServerResponse) (*PrivateNetworkServerModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if response == nil {
		diags.AddError("Constructor Error", "private network response is nil")
		return nil, diags
	}

	model := &PrivateNetworkServerModel{}

	model.Id = types.StringValue(response.Id)
	model.Name = types.StringValue(response.Name)
	if response.Lock != nil {
		model.Lock = types.Int64Value(*response.Lock)
	} else {
		model.Lock = types.Int64Null()
	}

	model.PrivateNetworkId = types.StringValue(privateNetworkId)

	return model, diags
}

func privateNetworkServerObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":   types.StringType,
			"name": types.StringType,
			"lock": types.Int64Type,
		},
	}
}

func PrivateNetworkServerDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for private network server information",
		Attributes: map[string]schema.Attribute{
			"private_network_id": schema.StringAttribute{
				Required:    true,
				Description: "Private network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid Id (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Private network server identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid Id (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Private network server name",
			},
			"lock": schema.Int64Attribute{
				Computed:    true,
				Description: "Private network server lock status",
			},
		},
	}
}
