package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicIpsModel struct {
	ID        types.String `tfsdk:"id"`
	PublicIps types.List   `tfsdk:"public_ips"`
}

func publicIpObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":            types.StringType,
			"ip":            types.StringType,
			"type":          types.StringType,
			"assigned_to":   assignedToObjectType(),
			"subnet_id":     types.StringType,
			"reverse_dns":   types.StringType,
			"is_dhcp":       types.BoolType,
			"state":         types.StringType,
			"datacenter":    baseDatacenterObjectType(),
			"creation_date": types.StringType,
		},
	}
}

func publicIpNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := PublicIpDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Public IP identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func NewPublicIps(ctx context.Context, publicIpsResponse []PublicIpResponse) (*PublicIpsModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &PublicIpsModel{}
	model.ID = types.StringValue("public_ips")

	publicIpModels, listDiags := NewPublicIpFromList(ctx, publicIpsResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		publicIpsList, convertDiags := types.ListValueFrom(ctx, publicIpObjectType(), publicIpModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.PublicIps = publicIpsList
		}
	}

	return model, diags
}

func PublicIpsDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing public IPs",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"public_ips": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of public IPs",
				NestedObject: publicIpNestedAttributeObject(),
			},
		},
	}
}
