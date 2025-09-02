package models

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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
	NetsSameVlan       types.Int64  `tfsdk:"nets_same_vlan"`
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
	NetsSameVlan       int     `json:"nets_same_vlan"`
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
	model.NetsSameVlan = types.Int64Value(int64(ip.NetsSameVlan))
	model.Type = types.StringValue(ip.Type)
	model.State = types.StringValue(ip.State)

	return model, diags
}

func NewPublicNetworkIpModel(ctx context.Context, ip *PublicNetworkIpResponse) (*PublicNetworkIpModel, diag.Diagnostics) {
	return newPublicNetworkIpFromResponse(ctx, ip)
}

func NewPublicNetworkIpFromList(ctx context.Context, sshList []PublicNetworkIpResponse) ([]PublicNetworkIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []PublicNetworkIpModel

	if len(sshList) == 0 {
		return []PublicNetworkIpModel{}, diags
	}

	for i, ssh := range sshList {
		model, modelDiags := NewPublicNetworkIpModel(ctx, &ssh)
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
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Public network IP identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
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
			"nets_same_vlan": schema.Int64Attribute{
				Computed:    true,
				Description: "Networks in same VLAN",
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

//type SshKeyCreateRequest struct {
//	Name        string `json:"name"`
//	Description string `json:"description,omitempty"`
//	PublicKey   string `json:"public_key,omitempty"`
//}
//
//func (m *SshKeyModel) ToCreateRequest() SshKeyCreateRequest {
//	return SshKeyCreateRequest{
//		Name:        m.Name.ValueString(),
//		Description: m.Description.ValueString(),
//		PublicKey:   m.PublicKey.ValueString(),
//	}
//}
//
//func NewSshKeyResourceModel(ctx context.Context, ssh *SshKeyResponse) (*SshKeyModel, diag.Diagnostics) {
//	baseModel, diags := newSshKeyFromResponse(ctx, ssh)
//	if diags.HasError() {
//		return nil, diags
//	}
//
//	return baseModel, diags
//}
//
//func SshKeyResourceSchema(_ context.Context) rschema.Schema {
//	return rschema.Schema{
//		Description: "SSH key resource",
//		Attributes: map[string]rschema.Attribute{
//			"id": rschema.StringAttribute{
//				Computed:    true,
//				Description: "SSH key identifier",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"name": rschema.StringAttribute{
//				Required:    true,
//				Description: "SSH key name",
//				Validators: []validator.String{
//					stringvalidator.LengthAtMost(util.MaxNameLength),
//					stringvalidator.LengthAtLeast(1),
//				},
//			},
//			"description": rschema.StringAttribute{
//				Computed:    true,
//				Optional:    true,
//				Description: "SSH key description",
//				Validators: []validator.String{
//					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
//				},
//			},
//			"state": rschema.StringAttribute{
//				Computed:    true,
//				Description: "Current state of the SSH key",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"servers": rschema.ListNestedAttribute{
//				Computed:     true,
//				Description:  "List of servers associated with the SSH key",
//				NestedObject: IdentifierResourceNestedObject(),
//				PlanModifiers: []planmodifier.List{
//					listplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"md5": rschema.StringAttribute{
//				Computed:    true,
//				Description: "MD5 hash of the SSH key",
//				Validators: []validator.String{
//					stringvalidator.RegexMatches(
//						regexp.MustCompile(util.HexID32Pattern),
//						"must be a valid MD5 hash (32 hexadecimal characters)",
//					),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"public_key": rschema.StringAttribute{
//				Computed:    true,
//				Optional:    true,
//				Description: "SSH public key content",
//				Validators: []validator.String{
//					stringvalidator.LengthAtMost(util.MaxNameLength),
//					stringvalidator.LengthAtLeast(1),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"creation_date": rschema.StringAttribute{
//				Computed:    true,
//				Description: "SSH key creation date in ISO 8601 format",
//				Validators: []validator.String{
//					stringvalidator.RegexMatches(
//						regexp.MustCompile(util.DateTimePattern),
//						"must be a date in ISO 8601 format (e.g., 2023-05-29T09:43:31+00:00)",
//					),
//				},
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//			"private_key": rschema.StringAttribute{
//				Computed:    true,
//				Description: "SSH key private key",
//				PlanModifiers: []planmodifier.String{
//					stringplanmodifier.UseStateForUnknown(),
//				},
//			},
//		},
//	}
//}
//
//type SshKeyUpdateRequest struct {
//	Name        string `json:"name"`
//	Description string `json:"description,omitempty"`
//}
//
//func (m *SshKeyModel) ToUpdateRequest() SshKeyUpdateRequest {
//	return SshKeyUpdateRequest{
//		Name:        m.Name.ValueString(),
//		Description: m.Description.ValueString(),
//	}
//}

//func (m *SshKeyModel) GetState() string {
//	return m.State.ValueString()
//}
