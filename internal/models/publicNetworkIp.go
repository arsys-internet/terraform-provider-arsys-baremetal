package models

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicNetworkIpModel struct {
	PublicNetworkId    types.String `tfsdk:"public_network_id"`
	Id                 types.String `tfsdk:"id"`
	IpAddress          types.String `tfsdk:"ip_address"`
	Description        types.String `tfsdk:"description"`
	NetworkInterfaceId types.String `tfsdk:"network_interface_id"`
	LbId               types.String `tfsdk:"lb_id"`
	InverseDns         types.String `tfsdk:"inverse_dns"`
	StartDate          types.String `tfsdk:"start_date"`
	SiteId             types.String `tfsdk:"site_id"`
	IsMain             types.Int64  `tfsdk:"is_main"`
	Mask               types.Int64  `tfsdk:"mask"`
	FirewallId         types.String `tfsdk:"firewall_id"`
	Gateway            types.String `tfsdk:"gateway"`
	Broadcast          types.String `tfsdk:"broadcast"`
	NetworkId          types.String `tfsdk:"network_id"`
	Type               types.String `tfsdk:"type"`
	State              types.String `tfsdk:"state"`
}

type PublicNetworkIpResponse struct {
	Id                 string  `json:"id"`
	IpAddress          string  `json:"ip_address"`
	Description        string  `json:"description"`
	NetworkInterfaceId *string `json:"network_interface_id"`
	LbId               *string `json:"lb_id"`
	InverseDns         string  `json:"inverse_dns"`
	StartDate          string  `json:"start_date"`
	SiteId             string  `json:"site_id"`
	IsMain             int     `json:"is_main"`
	Mask               int     `json:"mask"`
	FirewallId         *string `json:"firewall_id"`
	Gateway            *string `json:"gateway"`
	Broadcast          *string `json:"broadcast"`
	NetworkId          string  `json:"network_id"`
	Type               string  `json:"type"`
	State              string  `json:"state"`
}

func newPublicNetworkIpFromResponse(_ context.Context, ip *PublicNetworkIpResponse) (*PublicNetworkIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if ip == nil {
		diags.AddError("Response Error", "Public network ip response is nil")
		return nil, diags
	}

	model := &PublicNetworkIpModel{}

	model.PublicNetworkId = types.StringValue(ip.NetworkId)
	model.Id = types.StringValue(ip.Id)
	model.IpAddress = types.StringValue(ip.IpAddress)
	model.Description = types.StringValue(ip.Description)
	if ip.NetworkInterfaceId != nil {
		model.NetworkInterfaceId = types.StringValue(*ip.NetworkInterfaceId)
	} else {
		model.NetworkInterfaceId = types.StringNull()
	}
	if ip.LbId != nil {
		model.LbId = types.StringValue(*ip.LbId)
	} else {
		model.LbId = types.StringNull()
	}
	model.InverseDns = types.StringValue(ip.InverseDns)
	model.StartDate = types.StringValue(ip.StartDate)
	model.SiteId = types.StringValue(ip.SiteId)
	model.IsMain = types.Int64Value(int64(ip.IsMain))
	model.Mask = types.Int64Value(int64(ip.Mask))
	if ip.FirewallId != nil {
		model.FirewallId = types.StringValue(*ip.FirewallId)
	} else {
		model.FirewallId = types.StringNull()
	}
	if ip.Gateway != nil {
		model.Gateway = types.StringValue(*ip.Gateway)
	} else {
		model.Gateway = types.StringNull()
	}
	if ip.Broadcast != nil {
		model.Broadcast = types.StringValue(*ip.Broadcast)
	} else {
		model.Broadcast = types.StringNull()
	}
	model.NetworkId = types.StringValue(ip.NetworkId)
	model.Type = types.StringValue(ip.Type)
	model.State = types.StringValue(ip.State)

	return model, diags
}

func NewPublicNetworkIpModel(ctx context.Context, ip *PublicNetworkIpResponse) (*PublicNetworkIpModel, diag.Diagnostics) {
	return newPublicNetworkIpFromResponse(ctx, ip)
}

func NewPublicNetworkIpFromList(ctx context.Context, ipList []PublicNetworkIpResponse) ([]PublicNetworkIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []PublicNetworkIpModel

	if len(ipList) == 0 {
		return []PublicNetworkIpModel{}, diags
	}

	for i, ip := range ipList {
		model, modelDiags := NewPublicNetworkIpModel(ctx, &ip)
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

func PublicNetworkIpDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for public network IP information",
		Attributes: map[string]schema.Attribute{
			"public_network_id": schema.StringAttribute{
				Required:    true,
				Description: "Public network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID",
					),
				},
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Public network IP identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID",
					),
				},
			},
			"ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "IP address",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "IP description",
			},
			"network_interface_id": schema.StringAttribute{
				Computed:    true,
				Description: "Network interface identifier",
			},
			"lb_id": schema.StringAttribute{
				Computed:    true,
				Description: "Load balancer identifier",
			},
			"inverse_dns": schema.StringAttribute{
				Computed:    true,
				Description: "Inverse DNS name",
			},
			"start_date": schema.StringAttribute{
				Computed:    true,
				Description: "IP creation date",
			},
			"site_id": schema.StringAttribute{
				Computed:    true,
				Description: "Site identifier",
			},
			"is_main": schema.Int64Attribute{
				Computed:    true,
				Description: "IP is main",
			},
			"mask": schema.Int64Attribute{
				Computed:    true,
				Description: "IP mask",
			},
			"firewall_id": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall identifier",
			},
			"gateway": schema.StringAttribute{
				Computed:    true,
				Description: "Gateway address",
			},
			"broadcast": schema.StringAttribute{
				Computed:    true,
				Description: "Broadcast address",
			},
			"network_id": schema.StringAttribute{
				Computed:    true,
				Description: "Network identifier",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "IP type",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "IP state",
			},
		},
	}
}

type PublicNetworkIpResourceModel struct {
	PublicNetworkId types.String   `tfsdk:"public_network_id"`
	Id              types.String   `tfsdk:"id"`
	Action          types.Bool     `tfsdk:"action"`
	Ips             []types.String `tfsdk:"ips"`
	Items           types.List     `tfsdk:"items"`
}
type PublicNetworkIpRequest struct {
	PublicNetworkId string   `json:"public_network_id"`
	Action          bool     `json:"action"`
	Ips             []string `json:"ips"`
}

type PublicNetworkIpCreateResponse struct {
	Sync   bool    `json:"sync"`
	Data   IpsData `json:"data"`
	TaskId string  `json:"task_id"`
}

type IpsData struct {
	Items []PublicNetworkIpResponse `json:"items"`
}

func (m *PublicNetworkIpResourceModel) ToCreateRequest() PublicNetworkIpRequest {
	ips := make([]string, len(m.Ips))
	for i, ip := range m.Ips {
		ips[i] = ip.ValueString()
	}
	return PublicNetworkIpRequest{
		Action: m.Action.ValueBool(),
		Ips:    ips,
	}
}

func NewPublicNetworkIpResourceModel(ctx context.Context, data *PublicNetworkIpResourceModel, apiResponse []PublicNetworkIpResponse) (*PublicNetworkIpResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	ipModels, listDiags := NewPublicNetworkIpFromList(ctx, apiResponse)
	diags.Append(listDiags...)

	if !listDiags.HasError() {
		itemsList, convertDiags := types.ListValueFrom(ctx, publicNetworkIpObjectType(), ipModels)
		diags.Append(convertDiags...)
		if !convertDiags.HasError() {
			data.Items = itemsList
		}
	}

	return data, diags
}

func PublicNetworkIpResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Public network IP resource",
		Attributes: map[string]rschema.Attribute{
			"public_network_id": rschema.StringAttribute{
				Required:    true,
				Description: "Public network identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID",
					),
				},
			},
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Public network IP identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID",
					),
				},
			},
			"action": rschema.BoolAttribute{
				Required:    true,
				Description: "Action to perform on the IPs",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.RequiresReplace(),
				},
			},
			"ips": rschema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
				Description: "List of IP IDs to attach to the network",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.RequiresReplace(),
				},
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"items": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "List of attached IPs",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"public_network_id": rschema.StringAttribute{
							Computed:    true,
							Description: "Public network identifier",
						},
						"id": rschema.StringAttribute{
							Computed:    true,
							Description: "IP identifier",
						},
						"ip_address": rschema.StringAttribute{
							Computed:    true,
							Description: "IP address",
						},
						"description": rschema.StringAttribute{
							Computed:    true,
							Description: "IP description",
						},
						"network_interface_id": rschema.StringAttribute{
							Computed:    true,
							Description: "Network interface ID",
						},
						"lb_id": rschema.StringAttribute{
							Computed:    true,
							Description: "Load balancer ID",
						},
						"inverse_dns": rschema.StringAttribute{
							Computed:    true,
							Description: "Inverse DNS",
						},
						"start_date": rschema.StringAttribute{
							Computed:    true,
							Description: "Start date",
						},
						"site_id": rschema.StringAttribute{
							Computed:    true,
							Description: "Site identifier",
						},
						"is_main": rschema.Int64Attribute{
							Computed:    true,
							Description: "Is main IP",
						},
						"mask": rschema.Int64Attribute{
							Computed:    true,
							Description: "Network mask",
						},
						"firewall_id": rschema.StringAttribute{
							Computed:    true,
							Description: "Firewall identifier",
						},
						"gateway": rschema.StringAttribute{
							Computed:    true,
							Description: "Gateway address",
						},
						"broadcast": rschema.StringAttribute{
							Computed:    true,
							Description: "Broadcast address",
						},
						"network_id": rschema.StringAttribute{
							Computed:    true,
							Description: "Network identifier",
						},
						"type": rschema.StringAttribute{
							Computed:    true,
							Description: "IP type",
						},
						"state": rschema.StringAttribute{
							Computed:    true,
							Description: "IP state",
						},
					},
				},
			},
		},
	}
}

func (m *PublicNetworkIpModel) GetState() string {
	return m.State.ValueString()
}
