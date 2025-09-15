package models

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
	helper "terraform-provider-arsys-baremetal/internal/util/helper"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetworkModel struct {
	Id             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	NetworkAddress types.String `tfsdk:"network_address"`
	SubnetMask     types.String `tfsdk:"subnet_mask"`
	State          types.String `tfsdk:"state"`
	Datacenter     types.Object `tfsdk:"datacenter"`
	CreationDate   types.String `tfsdk:"creation_date"`
	Servers        types.List   `tfsdk:"servers"`
	CloudPanelId   types.String `tfsdk:"cloudpanel_id"`
}

type PrivateNetworkResourceModel struct {
	PrivateNetworkModel
	DatacenterId types.String `tfsdk:"datacenter_id"`
}

type PrivateNetworkResponse struct {
	Id             string                 `json:"id"`
	Name           string                 `json:"name"`
	Description    *string                `json:"description"`
	NetworkAddress string                 `json:"network_address"`
	SubnetMask     string                 `json:"subnet_mask"`
	State          string                 `json:"state"`
	Datacenter     BaseDatacenterResponse `json:"datacenter"`
	CreationDate   string                 `json:"creation_date"`
	Servers        []IdentifierResponse   `json:"servers"`
	CloudPanelId   string                 `json:"cloudpanel_id"`
}

func (f *PrivateNetworkModel) GetState() string {
	if f == nil {
		return ""
	}

	return f.State.ValueString()
}

func newPrivateNetworkModelFromResponse(_ context.Context, pn *PrivateNetworkResponse) (*PrivateNetworkModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if pn == nil {
		diags.AddError("Constructor Error", "private network response is nil")
		return nil, diags
	}

	model := &PrivateNetworkModel{}

	model.Id = types.StringValue(pn.Id)
	model.Name = types.StringValue(pn.Name)
	model.NetworkAddress = types.StringValue(pn.NetworkAddress)
	model.SubnetMask = types.StringValue(pn.SubnetMask)
	model.State = types.StringValue(pn.State)
	model.CreationDate = types.StringValue(pn.CreationDate)
	model.CloudPanelId = types.StringValue(pn.CloudPanelId)

	if pn.Description != nil {
		model.Description = types.StringValue(*pn.Description)
	} else {
		model.Description = types.StringNull()
	}

	datacenterObj, dcDiags := NewBaseDatacenterObject(pn.Datacenter)
	diags.Append(dcDiags...)
	if !dcDiags.HasError() {
		model.Datacenter = datacenterObj
	}

	serversList, listDiags := NewIdentifierList(pn.Servers)
	diags.Append(listDiags...)
	if !listDiags.HasError() {
		model.Servers = serversList
	}

	return model, diags
}

func NewPrivateNetworkModel(ctx context.Context, pn *PrivateNetworkResponse) (*PrivateNetworkModel, diag.Diagnostics) {
	return newPrivateNetworkModelFromResponse(ctx, pn)
}

func NewPrivateNetworkResourceModel(ctx context.Context, pn *PrivateNetworkResponse) (*PrivateNetworkResourceModel, diag.Diagnostics) {
	baseModel, diags := newPrivateNetworkModelFromResponse(ctx, pn)
	if diags.HasError() {
		return nil, diags
	}

	resourceModel := &PrivateNetworkResourceModel{
		PrivateNetworkModel: *baseModel,
		DatacenterId:        types.StringValue(pn.Datacenter.Id),
	}

	return resourceModel, diags
}

func privateNetworkObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":              types.StringType,
			"name":            types.StringType,
			"description":     types.StringType,
			"network_address": types.StringType,
			"subnet_mask":     types.StringType,
			"state":           types.StringType,
			"datacenter":      baseDatacenterObjectType(),
			"creation_date":   types.StringType,
			"servers":         types.ListType{ElemType: IdentifierObjectType()},
			"cloudpanel_id":   types.StringType,
		},
	}
}

func privateNetworkNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := PrivateNetworkDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Private network identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func NewPrivateNetworkFromList(ctx context.Context, pnList []PrivateNetworkResponse) ([]PrivateNetworkModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []PrivateNetworkModel

	if len(pnList) == 0 {
		return []PrivateNetworkModel{}, diags
	}

	for i, pn := range pnList {
		model, modelDiags := NewPrivateNetworkModel(ctx, &pn)
		if modelDiags.HasError() {
			diags.AddError(
				"List Constructor Error",
				fmt.Sprintf("Failed to create model for item %d: %s", i, modelDiags.Errors()[0].Summary()),
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

func (m *PrivateNetworkResourceModel) ToCreateRequest() PrivateNetworkCreateRequest {
	request := PrivateNetworkCreateRequest{
		Name:         m.Name.ValueString(),
		DatacenterId: m.DatacenterId.ValueString(),
	}

	helper.AssignStringPtr(&request.Description, m.Description)
	helper.AssignStringPtr(&request.NetworkAddress, m.NetworkAddress)
	helper.AssignStringPtr(&request.SubnetMask, m.SubnetMask)

	return request
}

func (m *PrivateNetworkResourceModel) ToUpdateRequestFromState(state *PrivateNetworkResourceModel) PrivateNetworkUpdateRequest {
	request := PrivateNetworkUpdateRequest{}

	if !m.Name.Equal(state.Name) {
		helper.AssignStringPtr(&request.Name, m.Name)
	}

	if !m.Description.Equal(state.Description) {
		helper.AssignStringPtr(&request.Description, m.Description)
	}

	networkChanged := !m.NetworkAddress.Equal(state.NetworkAddress)
	subnetChanged := !m.SubnetMask.Equal(state.SubnetMask)

	if networkChanged || subnetChanged {
		if !m.NetworkAddress.IsNull() {
			helper.AssignStringPtr(&request.NetworkAddress, m.NetworkAddress)
		} else {
			helper.AssignStringPtr(&request.NetworkAddress, state.NetworkAddress)
		}

		if !m.SubnetMask.IsNull() {
			helper.AssignStringPtr(&request.SubnetMask, m.SubnetMask)
		} else {
			helper.AssignStringPtr(&request.SubnetMask, state.SubnetMask)
		}
	}

	return request
}

type PrivateNetworkCreateRequest struct {
	Name           string  `json:"name"`
	Description    *string `json:"description,omitempty"`
	NetworkAddress *string `json:"network_address,omitempty"`
	SubnetMask     *string `json:"subnet_mask,omitempty"`
	DatacenterId   string  `json:"datacenter_id"`
}

type PrivateNetworkUpdateRequest struct {
	Name           *string `json:"name,omitempty"`
	Description    *string `json:"description,omitempty"`
	NetworkAddress *string `json:"network_address,omitempty"`
	SubnetMask     *string `json:"subnet_mask,omitempty"`
}

func PrivateNetworkDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for private network information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Private network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid private network ID",
					),
				},
			},
			"cloudpanel_id": schema.StringAttribute{
				Computed:    true,
				Description: "CloudPanel identifier",
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Private network name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Private network description",
			},
			"datacenter": BaseDatacenterNestedAttribute(),
			"network_address": schema.StringAttribute{
				Computed:    true,
				Description: "Network address",
			},
			"subnet_mask": schema.StringAttribute{
				Computed:    true,
				Description: "Subnet mask",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Private network state",
			},
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
			},
			"servers": schema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers in the private network",
				NestedObject: IdentifierNestedObject(),
			},
		},
	}
}

func PrivateNetworkResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Private network resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Private network identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Private network name",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxNameLength),
					stringvalidator.LengthAtLeast(1),
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.NamePattern),
						"must contain only alphanumeric characters, spaces, hyphens, underscores, and dots",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Private network description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Datacenter identifier where the network will be created",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter ID",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_address": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Network address",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.IPv4Pattern),
						"must be a valid IPv4 address (e.g., 192.168.1.0)",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_mask": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Subnet mask",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.SubnetMaskPattern),
						"must be a valid subnet mask (e.g., 255.255.255.0)",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Private network state",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"datacenter": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server datacenter",
				Attributes:  BaseDatacenterResourceAttributes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed:    true,
				Description: "CloudPanel identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"servers": rschema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers in the private network",
				NestedObject: IdentifierResourceNestedObject(),
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
