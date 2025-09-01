package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
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
						"must be a valid id (e.g., 4EEAD5836CF43ACA502FD5B99BFF44EF)",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"server_id": rschema.StringAttribute{
				Required:    true,
				Description: "Private network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid server_id (e.g., 4EEAD5836CF43ACA502FD5B99BFF44EF)",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
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
			"servers": rschema.ListNestedAttribute{
				Computed:     true,
				Description:  "List of servers in the private network",
				NestedObject: IdentifierResourceNestedObject(),
			},
		},
	}
}
