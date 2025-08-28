package models

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IdentifierModel struct {
	Id   types.String `tfsdk:"id"`
	Name types.String `tfsdk:"name"`
}

type IdentifierResponse struct {
	Id   string `json:"id"`
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
	model := IdentifierModel{
		Id:   types.StringValue(identifier.Id),
		Name: types.StringValue(identifier.Name),
	}

	return types.ObjectValueFrom(context.Background(), IdentifierObjectType().AttrTypes, model)
}

func NewIdentifierList(responses []IdentifierResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(responses) == 0 {
		return types.ListValueFrom(context.Background(), IdentifierObjectType(), []IdentifierModel{})
	}

	var models []IdentifierModel
	for _, response := range responses {
		model := IdentifierModel{
			Id:   types.StringValue(response.Id),
			Name: types.StringValue(response.Name),
		}
		models = append(models, model)
	}

	list, listDiags := types.ListValueFrom(context.Background(), IdentifierObjectType(), models)
	diags.Append(listDiags...)
	return list, diags
}

func BaseIdentifierAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Resource Id",
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

func BaseIdentifierResourceAttributes() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed:    true,
			Description: "Identifier",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"name": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Name",
		},
	}
}

func IdentifierNestedObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: BaseIdentifierAttributes(),
	}
}
