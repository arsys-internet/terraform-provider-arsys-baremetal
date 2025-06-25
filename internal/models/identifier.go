package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentifierModel struct {
	ID   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type IdentifierResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func IdentifierAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
	}
}

func IdentifierObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: IdentifierAttributeTypes(),
	}
}

func NewIdentifierObject(identifier IdentifierResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		IdentifierAttributeTypes(),
		map[string]attr.Value{
			"id":   types.StringValue(identifier.ID),
			"name": types.StringValue(identifier.Name),
		},
	)
}

func NewIdentifierList(responses []IdentifierResponse) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var values []attr.Value

	for _, server := range responses {
		serverObj, identifierDiags := NewIdentifierObject(server)
		diags.Append(identifierDiags...)
		if !identifierDiags.HasError() {
			values = append(values, serverObj)
		}
	}

	return types.ListValue(IdentifierObjectType(), values)
}

func BaseIdentifierAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Resource ID",
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Resource name",
		},
	}
}

func IdentifierResourceNestedObject() rschema.NestedAttributeObject {
	return rschema.NestedAttributeObject{
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Identifier",
			},
			"name": rschema.StringAttribute{
				Computed:    true,
				Description: "Name",
			},
		},
	}
}

func IdentifierNestedObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: BaseIdentifierAttributes(),
	}
}
