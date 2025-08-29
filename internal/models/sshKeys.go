package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type SshKeysModel struct {
	Id      types.String `tfsdk:"id"`
	SshKeys types.List   `tfsdk:"ssh_keys"`
}

func sshKeyObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":            types.StringType,
			"name":          types.StringType,
			"description":   types.StringType,
			"state":         types.StringType,
			"servers":       types.ListType{ElemType: IdentifierObjectType()},
			"md5":           types.StringType,
			"public_key":    types.StringType,
			"creation_date": types.StringType,
		},
	}
}

func sshKeyNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := SshKeyDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "SSH key identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func NewSshKeys(ctx context.Context, sshKeysResponse []SshKeyResponse) (*SshKeysModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := &SshKeysModel{}
	model.Id = types.StringValue("ssh_keys")

	sshKeyModels, listDiags := NewSshKeyFromList(ctx, sshKeysResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		sshKeysList, convertDiags := types.ListValueFrom(ctx, sshKeyObjectType(), sshKeyModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			model.SshKeys = sshKeysList
		}
	}

	return model, diags
}

func SshKeysDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for listing public IPs",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Identifier for this data source",
			},
			"ssh_keys": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of SSH keys",
				NestedObject: sshKeyNestedAttributeObject(),
			},
		},
	}
}
