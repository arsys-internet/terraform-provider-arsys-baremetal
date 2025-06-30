package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
)

type DatacenterModel struct {
	ID          types.String `tfsdk:"id"`
	CountryCode types.String `tfsdk:"country_code"`
	Location    types.String `tfsdk:"location"`
	Default     types.Int64  `tfsdk:"default"`
}

type DatacenterResponse struct {
	ID          string `json:"id"`
	CountryCode string `json:"country_code"`
	Location    string `json:"location"`
	Default     int64  `json:"default"`
}

func datacenterAttributeTypes() map[string]attr.Type {
	baseTypes := baseDatacenterAttributeTypes()
	baseTypes["default"] = types.Int64Type
	return baseTypes
}

func datacenterObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: datacenterAttributeTypes(),
	}
}

func NewDatacenter(_ context.Context, dc *DatacenterResponse) (*DatacenterModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if dc == nil {
		diags.AddError("Constructor Error", "datacenter response is nil")
		return nil, diags
	}

	model := &DatacenterModel{}
	model.ID = types.StringValue(dc.ID)
	model.CountryCode = types.StringValue(dc.CountryCode)
	model.Location = types.StringValue(dc.Location)
	model.Default = types.Int64Value(dc.Default)

	return model, diags
}

func NewDatacenterFromList(ctx context.Context, dcList []DatacenterResponse) ([]DatacenterModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []DatacenterModel

	if len(dcList) == 0 {
		return []DatacenterModel{}, diags
	}

	for i, dc := range dcList {
		model, modelDiags := NewDatacenter(ctx, &dc)
		if modelDiags.HasError() {
			diags.AddError(
				"Creation error",
				fmt.Sprintf("failed to create model for item %d: %s", i, modelDiags.Errors()[0].Summary()),
			)
			continue
		}
		diags.Append(modelDiags...)
		if model != nil {
			models = append(models, *model)
		}
	}

	return models, diags
}

func extendedDatacenterAttributes() map[string]schema.Attribute {
	attributes := baseDatacenterAttributes()
	attributes["default"] = schema.Int64Attribute{
		Computed:    true,
		Description: "Default datacenter flag",
	}
	return attributes
}

func datacenterNestedAttributeObject() schema.NestedAttributeObject {
	return schema.NestedAttributeObject{
		Attributes: extendedDatacenterAttributes(),
	}
}

func DatacenterDataSourceSchema(_ context.Context) schema.Schema {
	attributes := extendedDatacenterAttributes()
	attributes["id"] = schema.StringAttribute{
		Required:    true,
		Description: "Datacenter identifier",
		Validators: []validator.String{
			stringvalidator.RegexMatches(
				regexp.MustCompile(util.HexID32Pattern),
				"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
			),
		},
	}

	return schema.Schema{
		Description: "Data source for datacenter information",
		Attributes:  attributes,
	}
}
