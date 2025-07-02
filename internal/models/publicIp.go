package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
)

// PublicIpModel representa la estructura de datos para el modelo Terraform
type PublicIpModel struct {
	ID           types.String `tfsdk:"id"`
	IP           types.String `tfsdk:"ip"`
	Type         types.String `tfsdk:"type"`
	AssignedTo   types.Object `tfsdk:"assigned_to"`
	SubnetID     types.String `tfsdk:"subnet_id"`
	ReverseDNS   types.String `tfsdk:"reverse_dns"`
	IsDHCP       types.Bool   `tfsdk:"is_dhcp"`
	State        types.String `tfsdk:"state"`
	Datacenter   types.Object `tfsdk:"datacenter"`
	CreationDate types.String `tfsdk:"creation_date"`
}

// PublicIpResponse representa la estructura de la respuesta de la API
type PublicIpResponse struct {
	ID           string                 `json:"id"`
	IP           string                 `json:"ip"`
	Type         string                 `json:"type"`
	AssignedTo   AssignedToResponse     `json:"assigned_to"`
	SubnetID     *string                `json:"subnet_id"`
	ReverseDNS   *string                `json:"reverse_dns"`
	IsDHCP       bool                   `json:"is_dhcp"`
	State        string                 `json:"state"`
	Datacenter   BaseDatacenterResponse `json:"datacenter"`
	CreationDate string                 `json:"creation_date"`
}

// AssignedToResponse representa la estructura de la entidad asignada a la IP
type AssignedToResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// assignedToAttributeTypes returns the attribute types for AssignedTo
func assignedToAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":   types.StringType,
		"name": types.StringType,
		"type": types.StringType,
	}
}

// assignedToObjectType returns the object type for AssignedTo
func assignedToObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: assignedToAttributeTypes(),
	}
}

// NewAssignedToObject creates a types.Object from an AssignedToResponse
func NewAssignedToObject(assignedTo AssignedToResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		assignedToAttributeTypes(),
		map[string]attr.Value{
			"id":   types.StringValue(assignedTo.ID),
			"name": types.StringValue(assignedTo.Name),
			"type": types.StringValue(assignedTo.Type),
		},
	)
}

// NewPublicIp creates a PublicIpModel from the API response
func NewPublicIp(_ context.Context, ip *PublicIpResponse) (*PublicIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if ip == nil {
		diags.AddError("Response Error", "public ip response is nil")
		return nil, diags
	}

	model := &PublicIpModel{}

	model.ID = types.StringValue(ip.ID)
	model.IP = types.StringValue(ip.IP)
	model.Type = types.StringValue(ip.Type)
	model.IsDHCP = types.BoolValue(ip.IsDHCP)
	model.State = types.StringValue(ip.State)
	model.CreationDate = types.StringValue(ip.CreationDate)

	if ip.SubnetID != nil {
		model.SubnetID = types.StringValue(*ip.SubnetID)
	} else {
		model.SubnetID = types.StringNull()
	}

	if ip.ReverseDNS != nil {
		model.ReverseDNS = types.StringValue(*ip.ReverseDNS)
	} else {
		model.ReverseDNS = types.StringNull()
	}

	assignedToObj, assignedToDiags := NewAssignedToObject(ip.AssignedTo)
	diags.Append(assignedToDiags...)
	if !assignedToDiags.HasError() {
		model.AssignedTo = assignedToObj
	}

	datacenterObj, dcDiags := NewBaseDatacenterObject(ip.Datacenter)
	diags.Append(dcDiags...)
	if !dcDiags.HasError() {
		model.Datacenter = datacenterObj
	}

	return model, diags
}

// NewPublicIpFromList creates a list of models from a list of responses
func NewPublicIpFromList(ctx context.Context, ipList []PublicIpResponse) ([]PublicIpModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []PublicIpModel

	if len(ipList) == 0 {
		return []PublicIpModel{}, diags
	}

	for i, ip := range ipList {
		model, modelDiags := NewPublicIp(ctx, &ip)
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

// assignedToAttributes returns schema attributes for AssignedTo
func assignedToAttributes() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Computed:    true,
			Description: "ID of the resource to which the IP is assigned",
			Validators: []validator.String{
				stringvalidator.RegexMatches(
					regexp.MustCompile(util.HexID32Pattern),
					"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
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

// AssignedToNestedAttribute returns a nested attribute for AssignedTo
func AssignedToNestedAttribute() schema.SingleNestedAttribute {
	return schema.SingleNestedAttribute{
		Computed:    true,
		Description: "Information about the resource to which the IP is assigned",
		Attributes:  assignedToAttributes(),
	}
}

// PublicIpDataSourceSchema returns the schema for the PublicIp data source
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
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
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
				Description: "IP type (IPV4)",
				Validators: []validator.String{
					stringvalidator.OneOf("IPV4", "IPV6"),
				},
			},
			"assigned_to": AssignedToNestedAttribute(),
			"subnet_id": schema.StringAttribute{
				Computed:    true,
				Description: "ID of the subnet to which the IP belongs",
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
