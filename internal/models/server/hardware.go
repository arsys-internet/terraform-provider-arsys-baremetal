package server

import (
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
	"terraform-provider-arsys-baremetal/internal/util/helper"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type HardwareResponse struct {
	RAM                 int           `json:"ram"`
	HDDs                []HDDResponse `json:"hdds"`
	FixedInstanceSizeId *string       `json:"fixed_instance_size_id"`
	BaremetalModelId    string        `json:"baremetal_model_id"`
	VCore               int           `json:"vcore"`
	CoresPerProcessor   int           `json:"cores_per_processor"`
}

type HDDResponse struct {
	Id            string `json:"id"`
	Size          int    `json:"size"`
	IsMain        bool   `json:"is_main"`
	DiskType      string `json:"disk_type"`
	DiskRaid      string `json:"disk_raid"`
	DiskRaidCount int    `json:"disk_raid_count"`
}

type HardwareCreateRequest struct {
	BaremetalModelId string `json:"baremetal_model_id"`
}

func HardwareCreateRequestFromModel(hardwareObj types.Object) HardwareCreateRequest {
	hardwareAttrs := hardwareObj.Attributes()

	return HardwareCreateRequest{
		BaremetalModelId: hardwareAttrs["baremetal_model_id"].(types.String).ValueString(),
	}
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

	var fixedInstanceSizeId types.String
	if hardware.FixedInstanceSizeId != nil {
		fixedInstanceSizeId = types.StringValue(*hardware.FixedInstanceSizeId)
	} else {
		fixedInstanceSizeId = types.StringNull()
	}

	hardwareObj, objDiags := types.ObjectValue(HardwareObjectType().AttrTypes,
		map[string]attr.Value{
			"ram":                    types.Int64Value(int64(hardware.RAM)),
			"hdds":                   hddsList,
			"fixed_instance_size_id": fixedInstanceSizeId,
			"baremetal_model_id":     types.StringValue(hardware.BaremetalModelId),
			"vcore":                  types.Int64Value(int64(hardware.VCore)),
			"cores_per_processor":    types.Int64Value(int64(hardware.CoresPerProcessor)),
		})
	diags.Append(objDiags...)

	return hardwareObj, diags
}

func NeedsHardwareUpdate(hardwareAttrs map[string]attr.Value, apiHardware HardwareResponse) bool {
	if fixedId, exists := hardwareAttrs["fixed_instance_size_id"]; exists {
		if helper.GetStringValue(fixedId) != helper.GetStringPtr(apiHardware.FixedInstanceSizeId) {
			return true
		}
	}

	if modelId, exists := hardwareAttrs["baremetal_model_id"]; exists {
		if helper.GetStringValue(modelId) != apiHardware.BaremetalModelId {
			return true
		}
	}

	if ram, exists := hardwareAttrs["ram"]; exists {
		if ram.(types.Int64).ValueInt64() != int64(apiHardware.RAM) {
			return true
		}
	}

	if cores, exists := hardwareAttrs["cores_per_processor"]; exists {
		if cores.(types.Int64).ValueInt64() != int64(apiHardware.CoresPerProcessor) {
			return true
		}
	}

	if hdds, exists := hardwareAttrs["hdds"]; exists {
		currentCount := len(hdds.(types.List).Elements())
		if currentCount != len(apiHardware.HDDs) {
			return true
		}
	}

	return false
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
		"id":              types.StringValue(hdd.Id),
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
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(util.HexID32Pattern),
					"must be a valid baremetal model ID"),
			},
			Description: "Baremetal model identifier",
		},
		"vcore": rschema.Int64Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Description: "Number of virtual cores",
		},
		"cores_per_processor": rschema.Int64Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Description: "Number of cores per processor",
		},
		"ram": rschema.Int64Attribute{
			Computed: true,
			PlanModifiers: []planmodifier.Int64{
				int64planmodifier.UseStateForUnknown(),
			},
			Description: "Amount of RAM in GB",
		},
		"fixed_instance_size_id": rschema.StringAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
			Description: "Fixed instance size identifier",
		},
		"hdds": rschema.ListNestedAttribute{
			Computed: true,
			PlanModifiers: []planmodifier.List{
				listplanmodifier.UseStateForUnknown(),
			},
			Description: "List of hard disk drives",
			NestedObject: rschema.NestedAttributeObject{
				Attributes: map[string]rschema.Attribute{
					"size": rschema.Int64Attribute{
						Computed: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
						Description: "HDD size in GB",
					},
					"is_main": rschema.BoolAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
						Description: "Whether this is the main HDD",
					},
					"id": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "HDD identifier",
					},
					"disk_type": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "Type of disk",
					},
					"disk_raid": rschema.StringAttribute{
						Computed: true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
						Description: "RAID configuration",
					},
					"disk_raid_count": rschema.Int64Attribute{
						Computed: true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
						Description: "Number of disks in RAID",
					},
				},
			},
		},
	}
}
