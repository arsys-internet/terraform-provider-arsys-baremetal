package models

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strconv"
	"terraform-provider-arsys-baremetal/internal/util"
)

type OSArchitecture struct {
	Value int64
}

func (osa *OSArchitecture) UnmarshalJSON(data []byte) error {
	var i int64
	if err := json.Unmarshal(data, &i); err == nil {
		osa.Value = i
		return nil
	}

	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		switch s {
		case "x86", "32":
			osa.Value = 32
		case "x64", "64":
			osa.Value = 64
		case "arm64":
			osa.Value = 64
		default:
			if parsed, parseErr := strconv.ParseInt(s, 10, 64); parseErr == nil {
				osa.Value = parsed
			} else {
				return fmt.Errorf("valor de os_architecture no reconocido: %s", s)
			}
		}
		return nil
	}

	return fmt.Errorf("os_architecture debe ser string o int")
}

type ServerApplianceModel struct {
	ID                      types.String `tfsdk:"id"`
	Name                    types.String `tfsdk:"name"`
	AvailableDatacenters    types.List   `tfsdk:"available_datacenters"`
	OsFamily                types.String `tfsdk:"os_family"`
	Os                      types.String `tfsdk:"os"`
	OsVersion               types.String `tfsdk:"os_version"`
	OsArchitecture          types.Int64  `tfsdk:"os_architecture"`
	OsImageType             types.String `tfsdk:"os_image_type"`
	Type                    types.String `tfsdk:"type"`
	ServerTypeCompatibility types.List   `tfsdk:"server_type_compatibility"`
	MinHddSize              types.Int64  `tfsdk:"min_hdd_size"`
	Licenses                types.List   `tfsdk:"licenses"`
	Version                 types.String `tfsdk:"version"`
	Categories              types.List   `tfsdk:"categories"`
}

type ServerApplianceResponse struct {
	Id                      string         `json:"id"`
	Name                    string         `json:"name"`
	AvailableDatacenters    []string       `json:"available_datacenters"`
	OsFamily                string         `json:"os_family"`
	Os                      string         `json:"os"`
	OsVersion               string         `json:"os_version"`
	OsArchitecture          OSArchitecture `json:"os_architecture"`
	OsImageType             string         `json:"os_image_type"`
	Type                    string         `json:"type"`
	ServerTypeCompatibility []string       `json:"server_type_compatibility"`
	MinHddSize              int            `json:"min_hdd_size"`
	Licenses                []string       `json:"licenses"`
	Version                 string         `json:"version"`
	Categories              []string       `json:"categories"`
}

func NewServerAppliance(_ context.Context, sa *ServerApplianceResponse) (*ServerApplianceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if sa == nil {
		diags.AddError("Constructor Error", "server appliance response is nil")
		return nil, diags
	}

	model := &ServerApplianceModel{}

	model.ID = types.StringValue(sa.Id)
	model.Name = types.StringValue(sa.Name)
	model.OsFamily = types.StringValue(sa.OsFamily)
	model.Os = types.StringValue(sa.Os)
	model.OsVersion = types.StringValue(sa.OsVersion)
	model.OsArchitecture = types.Int64Value(sa.OsArchitecture.Value)
	model.Type = types.StringValue(sa.Type)
	model.MinHddSize = types.Int64Value(int64(sa.MinHddSize))
	model.Version = types.StringValue(sa.Version)
	model.OsImageType = types.StringValue(sa.OsImageType)

	elements := make([]attr.Value, len(sa.AvailableDatacenters))
	for i, dc := range sa.AvailableDatacenters {
		elements[i] = types.StringValue(dc)
	}
	listValue, listDiags := types.ListValue(types.StringType, elements)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.AvailableDatacenters = listValue
	}

	elements = make([]attr.Value, len(sa.ServerTypeCompatibility))
	for i, compat := range sa.ServerTypeCompatibility {
		elements[i] = types.StringValue(compat)
	}
	listValue, listDiags = types.ListValue(types.StringType, elements)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.ServerTypeCompatibility = listValue
	}

	elements = make([]attr.Value, len(sa.Licenses))
	for i, license := range sa.Licenses {
		elements[i] = types.StringValue(license)
	}
	listValue, listDiags = types.ListValue(types.StringType, elements)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.Licenses = listValue
	}

	elements = make([]attr.Value, len(sa.Categories))
	for i, category := range sa.Categories {
		elements[i] = types.StringValue(category)
	}
	listValue, listDiags = types.ListValue(types.StringType, elements)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.Categories = listValue
	}

	return model, diags
}

func serverApplianceObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                        types.StringType,
			"name":                      types.StringType,
			"available_datacenters":     types.ListType{ElemType: types.StringType},
			"os_family":                 types.StringType,
			"os":                        types.StringType,
			"os_version":                types.StringType,
			"os_architecture":           types.Int64Type,
			"os_image_type":             types.StringType,
			"type":                      types.StringType,
			"server_type_compatibility": types.ListType{ElemType: types.StringType},
			"min_hdd_size":              types.Int64Type,
			"licenses":                  types.ListType{ElemType: types.StringType},
			"version":                   types.StringType,
			"categories":                types.ListType{ElemType: types.StringType},
		},
	}
}

func serverApplianceNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := ServerApplianceDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Server appliance identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func NewServerApplianceFromList(ctx context.Context, pnList []ServerApplianceResponse) ([]ServerApplianceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []ServerApplianceModel

	if len(pnList) == 0 {
		return []ServerApplianceModel{}, diags
	}

	for i, sa := range pnList {
		model, modelDiags := NewServerAppliance(ctx, &sa)
		if modelDiags.HasError() {
			diags.AddError(
				"List Constructor Error",
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

func ServerApplianceDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for server appliance information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Server appliance identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Server appliance name",
			},
			"available_datacenters": schema.ListAttribute{
				Computed:    true,
				Description: "List of datacenter IDs where the appliance is available",
				ElementType: types.StringType,
			},
			"os_family": schema.StringAttribute{
				Computed:    true,
				Description: "Operating system family (Windows, Linux, Other)",
				Validators: []validator.String{
					stringvalidator.OneOf("Windows", "Linux", "Other"),
				},
			},
			"os": schema.StringAttribute{
				Computed:    true,
				Description: "Operating system",
			},
			"os_version": schema.StringAttribute{
				Computed:    true,
				Description: "Operating system version",
			},
			"os_architecture": schema.Int64Attribute{
				Computed:    true,
				Description: "OS architecture (32 or 64 bit)",
			},
			"os_image_type": schema.StringAttribute{
				Computed:    true,
				Description: "OS image type (STANDARD, MINIMAL, ISO_OS, ISO_TOOL, or null)",
				Validators: []validator.String{
					stringvalidator.OneOf("STANDARD", "MINIMAL", "ISO_OS", "ISO_TOOL"),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Server appliance type (IMAGE, MY_IMAGE, APPLICATION, ISO)",
				Validators: []validator.String{
					stringvalidator.OneOf("IMAGE", "MY_IMAGE", "APPLICATION", "ISO"),
				},
			},
			"server_type_compatibility": schema.ListAttribute{
				Computed:    true,
				Description: "List of server types compatible with this appliance",
				ElementType: types.StringType,
			},
			"min_hdd_size": schema.Int64Attribute{
				Computed:    true,
				Description: "Minimum hard disk size required in GB",
			},
			"licenses": schema.ListAttribute{
				Computed:    true,
				Description: "List of license identifiers",
				ElementType: types.StringType,
			},
			"version": schema.StringAttribute{
				Computed:    true,
				Description: "Server appliance version",
			},
			"categories": schema.ListAttribute{
				Computed:    true,
				Description: "List of categories this appliance belongs to",
				ElementType: types.StringType,
			},
		},
	}
}
