package server

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HardwareResponse struct {
	RAM                 int           `json:"ram"`
	HDDs                []HDDResponse `json:"hdds"`
	FixedInstanceSizeID *string       `json:"fixed_instance_size_id,omitempty"`
	BaremetalModelID    *string       `json:"baremetal_model_id,omitempty"`
	VCore               int           `json:"vcore"`
	CoresPerProcessor   int           `json:"cores_per_processor"`
}

type HDDResponse struct {
	ID            string `json:"id"`
	Size          int    `json:"size"`
	IsMain        bool   `json:"is_main"`
	DiskType      string `json:"disk_type"`
	DiskRaid      string `json:"disk_raid"`
	DiskRaidCount int    `json:"disk_raid_count"`
}

type HardwareCreateRequest struct {
	VCore             int                `json:"vcore,omitempty"`
	CoresPerProcessor int                `json:"cores_per_processor,omitempty"`
	RAM               int                `json:"ram,omitempty"`
	BaremetalModelID  string             `json:"baremetal_model_id"`
	HDDs              []HDDCreateRequest `json:"hdds"`
}

type HDDCreateRequest struct {
	Size   int  `json:"size"`
	IsMain bool `json:"is_main"`
}

func NewHardwareObject(hardware HardwareResponse) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	elements := make([]attr.Value, 0, len(hardware.HDDs))
	for _, hdd := range hardware.HDDs {
		hddObj, objDiags := NewHDDObject(hdd)
		diags.Append(objDiags...)

		if !objDiags.HasError() {
			elements = append(elements, hddObj)
		}
	}

	hddsList, listDiags := types.ListValue(HDDObjectType(), elements)
	diags.Append(listDiags...)

	var fixedInstanceSizeID types.String
	if hardware.FixedInstanceSizeID != nil {
		fixedInstanceSizeID = types.StringValue(*hardware.FixedInstanceSizeID)
	} else {
		fixedInstanceSizeID = types.StringNull()
	}

	var baremetalModelID types.String
	if hardware.BaremetalModelID != nil {
		baremetalModelID = types.StringValue(*hardware.BaremetalModelID)
	} else {
		baremetalModelID = types.StringNull()
	}

	hardwareObj, objDiags := types.ObjectValue(HardwareObjectType().AttrTypes,
		map[string]attr.Value{
			"ram":                    types.Int64Value(int64(hardware.RAM)),
			"hdds":                   hddsList,
			"fixed_instance_size_id": fixedInstanceSizeID,
			"baremetal_model_id":     baremetalModelID,
			"vcore":                  types.Int64Value(int64(hardware.VCore)),
			"cores_per_processor":    types.Int64Value(int64(hardware.CoresPerProcessor)),
		})
	diags.Append(objDiags...)

	return hardwareObj, diags
}

func HardwareObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"ram":                    types.Int64Type,
			"hdds":                   types.ListType{ElemType: HDDObjectType()},
			"fixed_instance_size_id": types.StringType,
			"baremetal_model_id":     types.StringType,
			"vcore":                  types.Int64Type,
			"cores_per_processor":    types.Int64Type,
		},
	}
}

func NewHDDObject(hdd HDDResponse) (types.Object, diag.Diagnostics) {
	attrs := map[string]attr.Value{
		"id":              types.StringValue(hdd.ID),
		"size":            types.Int64Value(int64(hdd.Size)),
		"is_main":         types.BoolValue(hdd.IsMain),
		"disk_type":       types.StringValue(hdd.DiskType),
		"disk_raid":       types.StringValue(hdd.DiskRaid),
		"disk_raid_count": types.Int64Value(int64(hdd.DiskRaidCount)),
	}

	return types.ObjectValue(HDDObjectType().AttrTypes, attrs)
}

func HDDObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":              types.StringType,
			"size":            types.Int64Type,
			"is_main":         types.BoolType,
			"disk_type":       types.StringType,
			"disk_raid":       types.StringType,
			"disk_raid_count": types.Int64Type,
		},
	}
}

func HardwareDataSourceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"baremetal_model_id": schema.StringAttribute{
			Computed:    true,
			Description: "Baremetal model identifier",
		},
		"vcore": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of virtual cores",
		},
		"cores_per_processor": schema.Int64Attribute{
			Computed:    true,
			Description: "Number of cores per processor",
		},
		"ram": schema.Int64Attribute{
			Computed:    true,
			Description: "Amount of RAM in GB",
		},
		"fixed_instance_size_id": schema.StringAttribute{
			Computed:    true,
			Description: "Fixed instance size identifier",
		},
		"hdds": schema.ListNestedAttribute{
			Computed:    true,
			Description: "List of hard disk drives",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Computed:    true,
						Description: "HDD identifier",
					},
					"size": schema.Int64Attribute{
						Computed:    true,
						Description: "HDD size in GB",
					},
					"is_main": schema.BoolAttribute{
						Computed:    true,
						Description: "Whether this is the main HDD",
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
				},
			},
		},
	}
}

func HardwareResourceSchema() map[string]rschema.Attribute {
	return map[string]rschema.Attribute{
		"baremetal_model_id": rschema.StringAttribute{
			Required:    true,
			Description: "Baremetal model identifier",
		},
		"vcore": rschema.Int64Attribute{
			Optional:    true,
			Description: "Number of virtual cores",
		},
		"cores_per_processor": rschema.Int64Attribute{
			Optional:    true,
			Description: "Number of cores per processor",
		},
		"ram": rschema.Int64Attribute{
			Optional:    true,
			Description: "Amount of RAM in GB",
		},
		"fixed_instance_size_id": rschema.StringAttribute{
			Computed:    true,
			Description: "Fixed instance size identifier",
		},
		"hdds": rschema.ListNestedAttribute{
			Required:    true,
			Description: "List of hard disk drives",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"size": rschema.Int64Attribute{
						Required:    true,
						Description: "HDD size in GB",
					},
					"is_main": rschema.BoolAttribute{
						Optional:    true,
						Computed:    true,
						Default:     booldefault.StaticBool(false),
						Description: "Whether this is the main HDD",
					},
					// Computed fields
					"id": rschema.StringAttribute{
						Computed:    true,
						Description: "HDD identifier",
					},
					"disk_type": rschema.StringAttribute{
						Computed:    true,
						Description: "Type of disk",
					},
					"disk_raid": rschema.StringAttribute{
						Computed:    true,
						Description: "RAID configuration",
					},
					"disk_raid_count": rschema.Int64Attribute{
						Computed:    true,
						Description: "Number of disks in RAID",
					},
				},
			},
		},
	}
}
