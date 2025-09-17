package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetworkServerRemoveModel struct {
	Id             types.String `tfsdk:"id"`
	ServerId       types.String `tfsdk:"server_id"`
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

func NewPrivateNetworkServerRemoveModel(_ context.Context, serverId string, pn PrivateNetworkResponse) (*PrivateNetworkServerRemoveModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	serversDetailList, listDiags := NewIdentifierList(pn.Servers)
	diags.Append(listDiags...)

	var description types.String
	if pn.Description != nil {
		description = types.StringValue(*pn.Description)
	} else {
		description = types.StringNull()
	}

	datacenterObj, dcDiags := NewBaseDatacenterObject(pn.Datacenter)
	diags.Append(dcDiags...)

	model := &PrivateNetworkServerRemoveModel{
		Id:           types.StringValue(pn.Id),
		ServerId:     types.StringValue(serverId),
		Name:         types.StringValue(pn.Name),
		Description:  description,
		State:        types.StringValue(pn.State),
		CreationDate: types.StringValue(pn.CreationDate),
		Datacenter:   datacenterObj,
		Servers:      serversDetailList,
		CloudPanelId: types.StringValue(pn.CloudPanelId),
	}

	return model, diags
}

func (m *PrivateNetworkServerRemoveModel) ToAssignRequest(ctx context.Context) (*PrivateNetworkServerRequest, diag.Diagnostics) {
	var servers []string

	diags := m.Servers.ElementsAs(ctx, &servers, false)
	if diags.HasError() {
		return nil, diags
	}

	return &PrivateNetworkServerRequest{
		Servers: servers,
	}, nil
}

func PrivateNetworkServerResourceRemoveSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Private network resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Required:    true,
				Description: "Private network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid private network ID",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_id": rschema.StringAttribute{
				Required:    true,
				Description: "Server identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid server ID",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": rschema.StringAttribute{
				Computed:    true,
				Description: "Private network name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"description": rschema.StringAttribute{
				Computed:    true,
				Description: "Private network description",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_address": rschema.StringAttribute{
				Computed:    true,
				Description: "Network address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_mask": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet mask",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Private network state",
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
			},
		},
	}
}
