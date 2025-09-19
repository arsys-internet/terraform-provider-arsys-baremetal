package models

import (
	"context"
	"fmt"
	"regexp"
	firewallpolicy "terraform-provider-arsys-baremetal/internal/models/firewallpolicy"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type FirewallPolicyResponse struct {
	Id           string                                    `json:"id"`
	Name         string                                    `json:"name"`
	Description  *string                                   `json:"description"`
	State        string                                    `json:"state"`
	CreationDate string                                    `json:"creation_date"`
	Default      int                                       `json:"default"`
	Rules        []firewallpolicy.FirewallRuleResponse     `json:"rules"`
	ServerIPs    []firewallpolicy.FirewallServerIPResponse `json:"server_ips"`
	CloudPanelId string                                    `json:"cloudpanel_id"`
}

type FirewallPolicyModel struct {
	Id           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	Description  types.String `tfsdk:"description"`
	State        types.String `tfsdk:"state"`
	CreationDate types.String `tfsdk:"creation_date"`
	Default      types.Int64  `tfsdk:"default"`
	Rules        types.List   `tfsdk:"rules"`
	ServerIPs    types.List   `tfsdk:"server_ips"`
	CloudPanelId types.String `tfsdk:"cloudpanel_id"`
}

func NewFirewallPolicyModel(_ context.Context, fp FirewallPolicyResponse) (*FirewallPolicyModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	rulesList, rulesDiags := firewallpolicy.NewFirewallRulesList(fp.Rules)
	diags.Append(rulesDiags...)

	serverIPsList, serverIPsDiags := firewallpolicy.NewFirewallServerIPsList(fp.ServerIPs)
	diags.Append(serverIPsDiags...)

	var description types.String
	if fp.Description != nil {
		description = types.StringValue(*fp.Description)
	} else {
		description = types.StringNull()
	}

	model := &FirewallPolicyModel{
		Id:           types.StringValue(fp.Id),
		Name:         types.StringValue(fp.Name),
		Description:  description,
		State:        types.StringValue(fp.State),
		CreationDate: types.StringValue(fp.CreationDate),
		Default:      types.Int64Value(int64(fp.Default)),
		Rules:        rulesList,
		ServerIPs:    serverIPsList,
		CloudPanelId: types.StringValue(fp.CloudPanelId),
	}

	return model, diags
}

func NewFirewallPolicyModelFromRead(_ context.Context, fp *FirewallPolicyResponse, currentState *FirewallPolicyModel) (*FirewallPolicyModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if currentState == nil || currentState.State.IsNull() || currentState.Name.IsNull() {
		return NewFirewallPolicyModel(context.Background(), *fp)
	}

	model := *currentState

	if fp.State != currentState.State.ValueString() {
		model.State = types.StringValue(fp.State)
	}

	return &model, diags
}

func FirewallPolicyObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":            types.StringType,
			"name":          types.StringType,
			"description":   types.StringType,
			"state":         types.StringType,
			"creation_date": types.StringType,
			"default":       types.Int64Type,
			"rules":         types.ListType{ElemType: firewallpolicy.FirewallRuleObjectType()},
			"server_ips":    types.ListType{ElemType: firewallpolicy.FirewallServerIPObjectType()},
			"cloudpanel_id": types.StringType,
		},
	}
}

type FirewallPolicyCreateRequest struct {
	Name        string                                     `json:"name"`
	Description *string                                    `json:"description,omitempty"`
	Rules       []firewallpolicy.FirewallRuleCreateRequest `json:"rules"`
}

type FirewallPolicyUpdateRequest struct {
	Name        string  `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (f *FirewallPolicyModel) GetState() string {
	if f == nil {
		return ""
	}

	return f.State.ValueString()
}

func (f *FirewallPolicyModel) ToCreateRequest() (FirewallPolicyCreateRequest, error) {
	request := FirewallPolicyCreateRequest{
		Name: f.Name.ValueString(),
	}

	if !f.Description.IsNull() {
		desc := f.Description.ValueString()
		request.Description = &desc
	}

	rules, err := firewallpolicy.ConvertRulesToCreateRequest(f.Rules)
	if err != nil {
		return request, fmt.Errorf("failed to convert rules: %w", err)
	}
	request.Rules = rules

	return request, nil
}

func (f *FirewallPolicyModel) ToUpdateRequest() FirewallPolicyUpdateRequest {
	request := FirewallPolicyUpdateRequest{}

	if !f.Name.IsNull() && f.Name.ValueString() != "" {
		request.Name = f.Name.ValueString()
	}

	if !f.Description.IsNull() {
		desc := f.Description.ValueString()
		request.Description = &desc
	}

	return request
}

func FirewallPolicyDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Fetches information about a specific firewall policy",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Firewall policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid firewall policy ID",
					),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy description",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy state",
			},
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "Date when firewall policy was created",
			},
			"default": schema.Int64Attribute{
				Computed:    true,
				Description: "Define default panel firewalls",
			},
			"cloudpanel_id": schema.StringAttribute{
				Computed:    true,
				Description: "Public identifier shown in panel",
			},
			"rules": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Firewall policy rules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: firewallpolicy.FirewallRuleDataSourceSchema(),
				},
			},
			"server_ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Servers assigned to firewall policy",
				NestedObject: schema.NestedAttributeObject{
					Attributes: firewallpolicy.FirewallServerIPDataSourceSchema(),
				},
			},
		},
	}
}

func FirewallPolicyResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Manages a firewall policy",
		Attributes: map[string]rschema.Attribute{
			"name": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Firewall policy name",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.NamePattern),
						"must contain only alphanumeric characters, spaces, hyphens, underscores, and dots",
					),
				},
			},
			"description": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Firewall policy description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
				},
			},
			"rules": rschema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Firewall policy rules",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: firewallpolicy.FirewallRuleResourceSchema(),
				},
			},
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Firewall policy state",
			},
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "Date when firewall policy was created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default": rschema.Int64Attribute{
				Computed:    true,
				Description: "Define default panel firewalls",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Identifier of the cloud panel",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_ips": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "ServerIPs assigned to firewall policy",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: firewallpolicy.FirewallServerIPResourceSchema(),
				},
			},
		},
	}
}
