package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaseDatacenterModel struct {
	ID          types.String `tfsdk:"id"`
	CountryCode types.String `tfsdk:"country_code"`
	Location    types.String `tfsdk:"location"`
}

type BaseDatacenterResponse struct {
	ID          string `json:"id"`
	CountryCode string `json:"country_code"`
	Location    string `json:"location"`
}

func baseDatacenterAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":           types.StringType,
		"country_code": types.StringType,
		"location":     types.StringType,
	}
}

func baseDatacenterObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: baseDatacenterAttributeTypes(),
	}
}

func NewBaseDatacenterObject(datacenter BaseDatacenterResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		baseDatacenterAttributeTypes(),
		map[string]attr.Value{
			"id":           types.StringValue(datacenter.ID),
			"country_code": types.StringValue(datacenter.CountryCode),
			"location":     types.StringValue(datacenter.Location),
		},
	)
}

func baseDatacenterAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Datacenter identifier",
		},
		"country_code": schema.StringAttribute{
			Computed:    true,
			Description: "Datacenter country code",
		},
		"location": schema.StringAttribute{
			Computed:    true,
			Description: "Datacenter location",
		},
	}
}

func BaseDatacenterNestedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Datacenter information",
		Attributes:  baseDatacenterAttributes(),
	}
}
