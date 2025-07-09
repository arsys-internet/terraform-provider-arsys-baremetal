package models

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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

func BaseDatacenterResourceAttributes() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"id": rschema.StringAttribute{
			Computed:    true,
			Description: "Datacenter ID",
		},
		"country_code": rschema.StringAttribute{
			Computed:    true,
			Description: "Country code",
		},
		"location": rschema.StringAttribute{
			Computed:    true,
			Description: "Location",
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

func BaseDatacenterResourceNestedAttribute() rschema.SingleNestedAttribute {
	return rschema.SingleNestedAttribute{
		Computed:    true,
		Description: "Datacenter information",
		Attributes:  BaseDatacenterResourceAttributes(),
	}
}
