package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BaremetalHardwareResponse struct {
	Core              int                    `json:"core"`
	CoresPerProcessor int                    `json:"cores_per_processor"`
	Ram               int                    `json:"ram"`
	Unit              string                 `json:"unit"`
	Hdds              []BaremetalHddResponse `json:"hdds"`
}

type BaremetalHddResponse struct {
	Size          int    `json:"size"`
	Unit          string `json:"unit"`
	DiskType      string `json:"disk_type"`
	DiskRaid      string `json:"disk_raid"`
	DiskRaidCount int    `json:"disk_raid_count"`
	IsMain        bool   `json:"is_main"`
}

type BaremetalHardwareModel struct {
	Core              types.Int64  `tfsdk:"core"`
	CoresPerProcessor types.Int64  `tfsdk:"cores_per_processor"`
	Ram               types.Int64  `tfsdk:"ram"`
	Unit              types.String `tfsdk:"unit"`
	Hdds              types.List   `tfsdk:"hdds"`
}

type BaremetalHddModel struct {
	Size          types.Int64  `tfsdk:"size"`
	Unit          types.String `tfsdk:"unit"`
	DiskType      types.String `tfsdk:"disk_type"`
	DiskRaid      types.String `tfsdk:"disk_raid"`
	DiskRaidCount types.Int64  `tfsdk:"disk_raid_count"`
	IsMain        types.Bool   `tfsdk:"is_main"`
}

func NewBaremetalHardwareObject(hardware BaremetalHardwareResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]attr.Value, 0, len(hardware.Hdds))
	for _, hdd := range hardware.Hdds {
		hddObj, objDiags := NewBaremetalHDDObject(hdd)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			elements = append(elements, hddObj)
		}
	}

	hddsList, listDiags := types.ListValue(BaremetalHDDObjectType(), elements)
	diags.Append(listDiags...)

	hardwareObj, objDiags := types.ObjectValue(BaremetalHardwareObjectType().AttrTypes,
		map[string]attr.Value{
			"core":                types.Int64Value(int64(hardware.Core)),
			"cores_per_processor": types.Int64Value(int64(hardware.CoresPerProcessor)),
			"ram":                 types.Int64Value(int64(hardware.Ram)),
			"unit":                types.StringValue(hardware.Unit),
			"hdds":                hddsList,
		})
	diags.Append(objDiags...)

	return hardwareObj, diags
}

func NewBaremetalHDDObject(hdd BaremetalHddResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"size":            types.Int64Value(int64(hdd.Size)),
		"unit":            types.StringValue(hdd.Unit),
		"disk_type":       types.StringValue(hdd.DiskType),
		"disk_raid":       types.StringValue(hdd.DiskRaid),
		"disk_raid_count": types.Int64Value(int64(hdd.DiskRaidCount)),
		"is_main":         types.BoolValue(hdd.IsMain),
	}

	return types.ObjectValue(BaremetalHDDObjectType().AttrTypes, attrs)
}

func BaremetalHDDObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"size":            types.Int64Type,
			"unit":            types.StringType,
			"disk_type":       types.StringType,
			"disk_raid":       types.StringType,
			"disk_raid_count": types.Int64Type,
			"is_main":         types.BoolType,
		},
	}
}

func BaremetalHardwareObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"core":                types.Int64Type,
			"cores_per_processor": types.Int64Type,
			"ram":                 types.Int64Type,
			"unit":                types.StringType,
			"hdds":                types.ListType{ElemType: BaremetalHDDObjectType()},
		},
	}
}

func BaremetalHardwareDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"core": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of cores",
		},
		"cores_per_processor": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of cores per processor",
		},
		"ram": schema.Int64Attribute{
			Computed:    true,
			Description: "Amount of RAM in GB",
		},
		"unit": schema.StringAttribute{
			Computed:    true,
			Description: "Unit of measurement",
		},
		"hdds": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of hard disk drives",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"size": schema.Int64Attribute{
						Computed:    true,
						Description: "HDD size",
					},
					"unit": schema.StringAttribute{
						Computed:    true,
						Description: "Unit of measurement",
					},
					"disk_type": schema.StringAttribute{
						Computed:    true,
						Description: "Type of disk",
					},
					"disk_raid": schema.StringAttribute{
						Computed:    true,
						Description: "RAID configuration",
					},
					"disk_raid_count": schema.Int64Attribute{
						Computed:    true,
						Description: "Number of disks in RAID",
					},
					"is_main": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this is the main HDD",
					},
				},
			},
		},
	}
}
