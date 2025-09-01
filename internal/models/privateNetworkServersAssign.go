package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PrivateNetworkServerAssignModel struct {
	Id             types.String `tfsdk:"id"`
	Servers        types.Set    `tfsdk:"servers"`
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	NetworkAddress types.String `tfsdk:"network_address"`
	SubnetMask     types.String `tfsdk:"subnet_mask"`
	State          types.String `tfsdk:"state"`
	Datacenter     types.Object `tfsdk:"datacenter"`
	CreationDate   types.String `tfsdk:"creation_date"`
	ServersDetail  types.List   `tfsdk:"servers_detail"`
	CloudPanelId   types.String `tfsdk:"cloudpanel_id"`
}

func NewPrivateNetworkServerAssignModel(_ context.Context, pn PrivateNetworkResponse) (*PrivateNetworkServerAssignModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	serverIds := make([]string, len(pn.Servers))
	for i, server := range pn.Servers {
		serverIds[i] = server.Id
	}
	serversList, _ := types.SetValueFrom(context.Background(), types.StringType, serverIds)

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

	model := &PrivateNetworkServerAssignModel{
		Id:            types.StringValue(pn.Id),
		Servers:       serversList,
		Name:          types.StringValue(pn.Name),
		Description:   description,
		State:         types.StringValue(pn.State),
		CreationDate:  types.StringValue(pn.CreationDate),
		Datacenter:    datacenterObj,
		ServersDetail: serversDetailList,
		CloudPanelId:  types.StringValue(pn.CloudPanelId),
	}

	return model, diags
}

func (m *PrivateNetworkServerAssignModel) ToAssignRequest(ctx context.Context) (*PrivateNetworkServerRequest, diag.Diagnostics) {
	var servers []string

	diags := m.Servers.ElementsAs(ctx, &servers, false)
	if diags.HasError() {
		return nil, diags
	}

	return &PrivateNetworkServerRequest{
		Servers: servers,
	}, nil
}

func PrivateNetworkServerAssignResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Private network resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Required:    true,
				Description: "Private network identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid id (e.g., 4EEAD5836CF43ACA502FD5B99BFF44EF)",
					),
				},
			},
			"servers": rschema.SetAttribute{
				ElementType: types.StringType,
				Required:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.RequiresReplace(),
				},
				Description: "List of server IP Ids to assign to the private network",
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
			"servers_detail": rschema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers in the private network",
				NestedObject: IdentifierResourceNestedObject(),
			},
		},
	}
}
