package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicNetworkServerRequest struct {
	Servers []string `json:"servers"`
}

type PublicNetworkServerResourceModel struct {
	PublicNetworkId   types.String   `tfsdk:"public_network_id"`
	Servers           []types.String `tfsdk:"servers"`
	Id                types.String   `tfsdk:"id"`
	PublicName        types.String   `tfsdk:"public_name"`
	Description       types.String   `tfsdk:"description"`
	DatacenterId      types.String   `tfsdk:"datacenter_id"`
	StartDate         types.Int64    `tfsdk:"start_date"`
	SameVlan          types.Bool     `tfsdk:"same_vlan"`
	Type              types.String   `tfsdk:"type"`
	State             types.String   `tfsdk:"state"`
	ServersDetails    types.List     `tfsdk:"servers_details"`
	Ips               types.List     `tfsdk:"ips"`
	AvailabilityZones types.List     `tfsdk:"availability_zones"`
	LastLogs          types.List     `tfsdk:"last_logs"`
}

func NewPublicNetworkServerResourceModel(ctx context.Context, publicNetworkId string, servers []string, apiResponse *PublicNetworkResponse) (*PublicNetworkServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	fullModel, fullModelDiags := NewPublicNetworkModel(ctx, apiResponse)
	diags.Append(fullModelDiags...)
	if diags.HasError() {
		return nil, diags
	}

	serverList := make([]types.String, len(servers))
	for i, server := range servers {
		serverList[i] = types.StringValue(server)
	}

	model := &PublicNetworkServerResourceModel{
		PublicNetworkId:   types.StringValue(publicNetworkId),
		Servers:           serverList,
		Id:                fullModel.Id,
		PublicName:        fullModel.PublicName,
		Description:       fullModel.Description,
		DatacenterId:      fullModel.DatacenterId,
		StartDate:         fullModel.StartDate,
		SameVlan:          fullModel.SameVlan,
		Type:              fullModel.Type,
		State:             fullModel.State,
		ServersDetails:    fullModel.Servers,
		Ips:               fullModel.Ips,
		AvailabilityZones: fullModel.AvailabilityZones,
		LastLogs:          fullModel.LastLogs,
	}

	return model, diags
}

func PublicNetworkServerSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Public network resource",
		Attributes: map[string]rschema.Attribute{
			"public_network_id": rschema.StringAttribute{
				Required:    true,
				Description: "Public network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid public network ID",
					),
				},
			},
			"servers": rschema.ListAttribute{
				Required:    true,
				Description: "List of servers identifiers in the public network",
				ElementType: types.StringType,
			},
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Public network identifier",
			},
			"public_name": rschema.StringAttribute{
				Computed:    true,
				Description: "Public network name",
			},
			"description": rschema.StringAttribute{
				Computed:    true,
				Description: "Public network description",
			},
			"datacenter_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Datacenter identifier where the network is located",
			},
			"start_date": rschema.Int64Attribute{
				Computed:    true,
				Description: "Start date timestamp in milliseconds",
			},
			"same_vlan": rschema.BoolAttribute{
				Computed:    true,
				Description: "Whether the network uses the same VLAN",
			},
			"type": rschema.StringAttribute{
				Computed:    true,
				Description: "Network type",
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Public network state",
			},
			"servers_details": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "Detailed list of servers in the public network",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"id": rschema.StringAttribute{
							Computed:    true,
							Description: "Server identifier",
						},
						"vlan_id": rschema.Int64Attribute{
							Computed:    true,
							Description: "VLAN identifier",
						},
						"mac": rschema.StringAttribute{
							Computed:    true,
							Description: "MAC address",
						},
						"tagged": rschema.BoolAttribute{
							Computed:    true,
							Description: "Whether the server is tagged",
						},
					},
				},
			},
			"ips": rschema.ListAttribute{
				Computed:    true,
				Description: "List of IP identifiers in the public network",
				ElementType: types.StringType,
			},
			"availability_zones": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "List of availability zones",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"id": rschema.StringAttribute{
							Computed:    true,
							Description: "Availability zone identifier",
						},
						"vlan_id": rschema.Int64Attribute{
							Computed:    true,
							Description: "VLAN identifier",
						},
					},
				},
			},
			"last_logs": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "List of recent log events",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"id": rschema.StringAttribute{
							Computed:    true,
							Description: "Log entry identifier",
						},
						"uuid": rschema.StringAttribute{
							Computed:    true,
							Description: "Log entry UUID",
						},
						"date": rschema.StringAttribute{
							Computed:    true,
							Description: "Log entry date",
						},
						"action": rschema.StringAttribute{
							Computed:    true,
							Description: "Action performed",
						},
						"time": rschema.Int64Attribute{
							Computed:    true,
							Description: "Time taken for the action",
						},
						"result": rschema.StringAttribute{
							Computed:    true,
							Description: "Result of the action",
						},
						"type": rschema.StringAttribute{
							Computed:    true,
							Description: "Log entry type",
						},
					},
				},
			},
		},
	}
}
