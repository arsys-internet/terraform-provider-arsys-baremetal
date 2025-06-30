package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
)

type PrivateNetworkModel struct {
	ID             types.String `tfsdk:"id"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	NetworkAddress types.String `tfsdk:"network_address"`
	SubnetMask     types.String `tfsdk:"subnet_mask"`
	State          types.String `tfsdk:"state"`
	Datacenter     types.Object `tfsdk:"datacenter"`
	CreationDate   types.String `tfsdk:"creation_date"`
	Servers        types.List   `tfsdk:"servers"`
	CloudPanelID   types.String `tfsdk:"cloudpanel_id"`
	DatacenterID   types.String `tfsdk:"datacenter_id"`
}

type PrivateNetworkResponse struct {
	ID             string                 `json:"id"`
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

func NewPrivateNetwork(_ context.Context, pn *PrivateNetworkResponse) (*PrivateNetworkModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if pn == nil {
		diags.AddError("Constructor Error", "private network response is nil")
		return nil, diags
	}

	model := &PrivateNetworkModel{}

	model.ID = types.StringValue(pn.ID)
	model.Name = types.StringValue(pn.Name)
	model.NetworkAddress = types.StringValue(pn.NetworkAddress)
	model.SubnetMask = types.StringValue(pn.SubnetMask)
	model.State = types.StringValue(pn.State)
	model.CreationDate = types.StringValue(pn.CreationDate)
	model.CloudPanelID = types.StringValue(pn.CloudPanelId)
	model.DatacenterID = types.StringValue(pn.Datacenter.ID)

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
			"datacenter_id":   types.StringType,
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
		model, modelDiags := NewPrivateNetwork(ctx, &pn)
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

func (m *PrivateNetworkModel) ToCreateRequest() PrivateNetworkCreateRequest {
	return PrivateNetworkCreateRequest{
		Name:           m.Name.ValueString(),
		Description:    m.Description.ValueString(),
		NetworkAddress: m.NetworkAddress.ValueString(),
		SubnetMask:     m.SubnetMask.ValueString(),
		DatacenterID:   m.DatacenterID.ValueString(),
	}
}

func (m *PrivateNetworkModel) ToUpdateRequest() PrivateNetworkUpdateRequest {
	return PrivateNetworkUpdateRequest{
		Name:           m.Name.ValueString(),
		Description:    m.Description.ValueString(),
		NetworkAddress: m.NetworkAddress.ValueString(),
		SubnetMask:     m.SubnetMask.ValueString(),
	}
}

type PrivateNetworkCreateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	NetworkAddress string `json:"network_address"`
	SubnetMask     string `json:"subnet_mask"`
	DatacenterID   string `json:"datacenter_id"`
}

type PrivateNetworkUpdateRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	NetworkAddress string `json:"network_address"`
	SubnetMask     string `json:"subnet_mask"`
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
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"cloudpanel_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
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
			"datacenter_id": schema.StringAttribute{
				Computed:    true,
				Description: "Datacenter identifier",
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
			},
			"description": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Private network description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Datacenter identifier where the network will be created",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter_id (e.g., 4EEAD5836CF43ACA502FD5B99BFF44EF)",
					),
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
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Private network state",
			},
			"datacenter": BaseDatacenterNestedAttribute(),
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed:    true,
				Description: "CloudPanel identifier",
			},
			"servers": rschema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers in the private network",
				NestedObject: IdentifierResourceNestedObject(),
			},
		},
	}
}
