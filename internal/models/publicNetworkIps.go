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

type PublicNetworkIpsModel struct {
	PublicNetworkId  types.String `tfsdk:"public_network_id"`
	Id               types.String `tfsdk:"id"`
	PublicNetworkIps types.List   `tfsdk:"ips_detail"`
}

func publicNetworkIpObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                   types.StringType,
			"ip_address":           types.StringType,
			"description":          types.StringType,
			"network_interface_id": types.StringType,
			"lb_id":                types.StringType,
			"inverse_dns":          types.StringType,
			"start_date":           types.StringType,
			"site_id":              types.StringType,
			"is_main":              types.Int64Type,
			"mask":                 types.Int64Type,
			"firewall_id":          types.StringType,
			"gateway":              types.StringType,
			"broadcast":            types.StringType,
			"network_id":           types.StringType,
			"nets_same_vlan":       types.Int64Type,
			"type":                 types.StringType,
			"state":                types.StringType,
		},
	}
}

func publicNetworkIpNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := PublicNetworkIpDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Public Network IP identifier",
			}
		} else if name == "public_network_id" {
			//
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func NewPublicNetworkIps(ctx context.Context, publicNetworkIpsResponse []PublicNetworkIpResponse) (*PublicNetworkIpsModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &PublicNetworkIpsModel{}
	model.Id = types.StringValue("public_network_ips")

	sshKeyModels, listDiags := NewPublicNetworkIpFromList(ctx, publicNetworkIpsResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		publicNetworkIpsList, convertDiags := types.ListValueFrom(ctx, publicNetworkIpObjectType(), sshKeyModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.PublicNetworkIps = publicNetworkIpsList
		}
	}

	return model, diags
}

func PublicNetworkIpsDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing IPs in the public network",
		Attributes: map[string]schema.Attribute{
			"public_network_id": schema.StringAttribute{
				Required:    true,
				Description: "Public network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"ips_detail": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of IPs in the public network",
				NestedObject: publicNetworkIpNestedAttributeObject(),
			},
		},
	}
}
