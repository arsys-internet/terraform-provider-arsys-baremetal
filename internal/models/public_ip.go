package models

import (
	"context"
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicIpModel struct {
	Id           types.String `tfsdk:"id"`
	IP           types.String `tfsdk:"ip"`
	Type         types.String `tfsdk:"type"`
	AssignedTo   types.Object `tfsdk:"assigned_to"`
	SubnetId     types.String `tfsdk:"subnet_id"`
	ReverseDNS   types.String `tfsdk:"reverse_dns"`
	IsDHCP       types.Bool   `tfsdk:"is_dhcp"`
	State        types.String `tfsdk:"state"`
	Datacenter   types.Object `tfsdk:"datacenter"`
	CreationDate types.String `tfsdk:"creation_date"`
}

type PublicIpResourceModel struct {
	PublicIpModel
	DatacenterId types.String `tfsdk:"datacenter_id"`
}

type PublicIpResponse struct {
	Id           string                 `json:"id"`
	IP           string                 `json:"ip"`
	Type         string                 `json:"type"`
	AssignedTo   *AssignedToResponse    `json:"assigned_to"`
	SubnetId     *string                `json:"subnet_id"`
	ReverseDNS   *string                `json:"reverse_dns"`
	IsDHCP       bool                   `json:"is_dhcp"`
	State        string                 `json:"state"`
	Datacenter   BaseDatacenterResponse `json:"datacenter"`
	CreationDate string                 `json:"creation_date"`
}

type AssignedToResponse struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func assignedToAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	}
}

func assignedToObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: assignedToAttributeTypes(),
	}
}

func NewAssignedToObject(assignedTo AssignedToResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		assignedToAttributeTypes(),
		map[string]attr.Value{
			"id":   types.StringValue(assignedTo.Id),
			"name": types.StringValue(assignedTo.Name),
			"type": types.StringValue(assignedTo.Type),
		},
	)
}

func newPublicIpFromResponse(_ context.Context, ip *PublicIpResponse) (*PublicIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if ip == nil {
		diags.AddError("Response Error", "public ip response is nil")
		return nil, diags
	}

	model := &PublicIpModel{}

	model.Id = types.StringValue(ip.Id)
	model.IP = types.StringValue(ip.IP)
	model.Type = types.StringValue(ip.Type)
	model.IsDHCP = types.BoolValue(ip.IsDHCP)
	model.State = types.StringValue(ip.State)
	model.CreationDate = types.StringValue(ip.CreationDate)

	if ip.SubnetId != nil {
		model.SubnetId = types.StringValue(*ip.SubnetId)
	} else {
		model.SubnetId = types.StringNull()
	}

	if ip.ReverseDNS != nil {
		model.ReverseDNS = types.StringValue(*ip.ReverseDNS)
	} else {
		model.ReverseDNS = types.StringNull()
	}

	if ip.AssignedTo != nil {
		assignedToObj, assignedToDiags := NewAssignedToObject(*ip.AssignedTo)
		diags.Append(assignedToDiags...)
		if !assignedToDiags.HasError() {
			model.AssignedTo = assignedToObj
		}
	} else {
		model.AssignedTo = types.ObjectNull(assignedToAttributeTypes())
	}

	datacenterObj, dcDiags := NewBaseDatacenterObject(ip.Datacenter)
	diags.Append(dcDiags...)
	if !dcDiags.HasError() {
		model.Datacenter = datacenterObj
	}

	return model, diags
}

func NewPublicIpModel(ctx context.Context, ip *PublicIpResponse) (*PublicIpModel, diag.Diagnostics) {
	return newPublicIpFromResponse(ctx, ip)
}

func NewPublicIpResourceModel(ctx context.Context, ip *PublicIpResponse) (*PublicIpResourceModel, diag.Diagnostics) {
	baseModel, diags := newPublicIpFromResponse(ctx, ip)
	if diags.HasError() {
		return nil, diags
	}

	resourceModel := &PublicIpResourceModel{
		PublicIpModel: *baseModel,
		DatacenterId:  types.StringValue(ip.Datacenter.Id),
	}

	return resourceModel, diags
}

func NewPublicIpFromList(ctx context.Context, ipList []PublicIpResponse) ([]PublicIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []PublicIpModel

	if len(ipList) == 0 {
		return []PublicIpModel{}, diags
	}

	for i, ip := range ipList {
		model, modelDiags := NewPublicIpModel(ctx, &ip)
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

func assignedToAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "Id of the resource to which the IP is assigned",
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(util.HexID32Pattern),
					"must be a valid ID",
				),
			},
		},
		"name": schema.StringAttribute{
			Computed:    true,
			Description: "Name of the resource to which the IP is assigned",
		},
		"type": schema.StringAttribute{
			Computed:    true,
			Description: "Type of resource to which the IP is assigned (SERVER, etc.)",
			Validators: []validator.String{
				stringvalidator.OneOf("SERVER", "LOAD_BALANCER"),
			},
		},
	}
}

func AssignedToNestedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Information about the resource to which the IP is assigned",
		Attributes:  assignedToAttributes(),
	}
}

func PublicIpDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for public IP information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Public IP identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid public IP ID",
					),
				},
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Description: "IP address",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.IPv4Pattern),
						"must be a valid IPv4 address",
					),
				},
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "IP type",
				Validators: []validator.String{
					stringvalidator.OneOf("IPV4", "IPV6"),
				},
			},
			"assigned_to": AssignedToNestedAttribute(),
			"subnet_id": schema.StringAttribute{
				Computed:    true,
				Description: "Id of the subnet to which the IP belongs",
			},
			"reverse_dns": schema.StringAttribute{
				Computed:    true,
				Description: "Reverse DNS configured for the IP",
			},
			"is_dhcp": schema.BoolAttribute{
				Computed:    true,
				Description: "Indicates if the IP is configured to use DHCP",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Current state of the IP (ACTIVE, etc.)",
			},
			"datacenter": BaseDatacenterNestedAttribute(),
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "IP creation date in ISO 8601 format",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.DateTimePattern),
						"must be a date in ISO 8601 format (e.g., 2023-05-29T09:43:31+00:00)",
					),
				},
			},
		},
	}
}

type PublicIpCreateRequest struct {
	ReverseDns   string `json:"reverse_dns,omitempty"`
	DatacenterId string `json:"datacenter_id,omitempty"`
	Type         string `json:"type,omitempty"`
}

func (m *PublicIpResourceModel) ToCreateRequest() PublicIpCreateRequest {
	return PublicIpCreateRequest{
		ReverseDns:   m.ReverseDNS.ValueString(),
		DatacenterId: m.DatacenterId.ValueString(),
		Type:         m.Type.ValueString(),
	}
}

func PublicIpResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Public ip resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Public ip identifier",
			},
			"reverse_dns": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Reverse DNS name",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxNameLength),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Datacenter identifier where the ip will be created",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter ID",
					),
				},
			},
			"type": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf("IPV4", "IPV6"),
				},
			},
			"ip": rschema.StringAttribute{
				Computed:    true,
				Description: "IP address",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"datacenter":  BaseDatacenterNestedAttribute(),
			"assigned_to": AssignedToNestedAttribute(),
			"subnet_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Id of the subnet to which the IP belongs",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"is_dhcp": rschema.BoolAttribute{
				Computed:    true,
				Description: "IP use DHCP",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "Current state of the IP (ACTIVE, etc.)",
			},
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "IP creation date",
			},
		},
	}
}

type PublicIpUpdateRequest struct {
	ReverseDns string `json:"reverse_dns"`
}

func (m *PublicIpResourceModel) ToUpdateRequest() PublicIpUpdateRequest {
	return PublicIpUpdateRequest{
		ReverseDns: m.ReverseDNS.ValueString(),
	}
}

func (m *PublicIpResourceModel) GetState() string {
	return m.State.ValueString()
}
