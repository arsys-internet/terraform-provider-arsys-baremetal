package models

import (
	"context"
	"regexp"
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

type SubnetModel struct {
	Id                     types.String `tfsdk:"id"`
	Name                   types.String `tfsdk:"name"`
	PublicName             types.String `tfsdk:"public_name"`
	Ip                     types.String `tfsdk:"ip"`
	Description            types.String `tfsdk:"description"`
	TypeId                 types.Int64  `tfsdk:"type_id"`
	Type                   types.String `tfsdk:"type"`
	Subnet                 types.String `tfsdk:"subnet"`
	SubnetId               types.String `tfsdk:"subnet_id"`
	NetworkInterfaceId     types.String `tfsdk:"network_interface_id"`
	ServerType             types.String `tfsdk:"server_type"`
	LoadBalancerId         types.String `tfsdk:"load_balancer_id"`
	LoadBalancerPublicName types.String `tfsdk:"load_balancer_public_name"`
	NetworkId              types.String `tfsdk:"network_id"`
	NetworkPublicName      types.String `tfsdk:"network_public_name"`
	Mask                   types.Int64  `tfsdk:"mask"`
	Gateway                types.String `tfsdk:"gateway"`
	Broadcast              types.String `tfsdk:"broadcast"`
	InverseDns             types.String `tfsdk:"inversedns"`
	Dhcp                   types.Int64  `tfsdk:"dhcp"`
	StateId                types.Int64  `tfsdk:"state_id"`
	State                  types.String `tfsdk:"state"`
	DatacenterId           types.String `tfsdk:"datacenter_id"`
	StartDate              types.String `tfsdk:"start_date"`
	LastLogs               types.List   `tfsdk:"last_logs"`
}

type SubnetResponse struct {
	Id                     string           `json:"id"`
	Name                   string           `json:"name"`
	PublicName             string           `json:"public_name"`
	Ip                     string           `json:"ip"`
	Description            *string          `json:"description"`
	TypeId                 int64            `json:"type_id"`
	Type                   string           `json:"type"`
	Subnet                 *string          `json:"subnet"`
	SubnetId               *string          `json:"subnet_id"`
	NetworkInterfaceId     *string          `json:"network_interface_id"`
	ServerType             *string          `json:"server_type"`
	LoadBalancerId         *string          `json:"load_balancer_id"`
	LoadBalancerPublicName *string          `json:"load_balancer_public_name"`
	NetworkId              *string          `json:"network_id"`
	NetworkPublicName      *string          `json:"network_public_name"`
	Mask                   int64            `json:"mask"`
	Gateway                string           `json:"gateway"`
	Broadcast              string           `json:"broadcast"`
	InverseDns             *string          `json:"inversedns"`
	Dhcp                   int64            `json:"dhcp"`
	StateId                int64            `json:"state_id"`
	State                  string           `json:"state"`
	DatacenterId           string           `json:"datacenter_id"`
	StartDate              string           `json:"start_date"`
	LastLogs               []SubnetLogEntry `json:"last_logs"`
}

type SubnetLogEntry struct {
	Id     string `json:"id"`
	Uuid   string `json:"uuid"`
	Date   string `json:"date"`
	Action string `json:"action"`
	Time   int64  `json:"time"`
	Result string `json:"result"`
	Type   string `json:"type"`
}

type CreateSubnetRequest struct {
	Mask         int64  `json:"mask"`
	DatacenterId string `json:"datacenter_id"`
}

func (s *SubnetModel) ToCreateRequest() CreateSubnetRequest {
	return CreateSubnetRequest{
		Mask:         s.Mask.ValueInt64(),
		DatacenterId: s.DatacenterId.ValueString(),
	}
}

func NewSubnetLogsList(logs []SubnetLogEntry) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(logs) == 0 {
		return types.ListNull(subnetLogEntryObjectType()), diags
	}

	elements := make([]attr.Value, len(logs))
	for i, log := range logs {
		obj, objDiags := types.ObjectValue(subnetLogEntryObjectType().AttrTypes, map[string]attr.Value{
			"id":     types.StringValue(log.Id),
			"uuid":   types.StringValue(log.Uuid),
			"date":   types.StringValue(log.Date),
			"action": types.StringValue(log.Action),
			"time":   types.Int64Value(log.Time),
			"result": types.StringValue(log.Result),
			"type":   types.StringValue(log.Type),
		})
		if objDiags.HasError() {
			diags.Append(objDiags...)
			continue
		}
		elements[i] = obj
	}

	if diags.HasError() {
		return types.ListNull(subnetLogEntryObjectType()), diags
	}

	return types.ListValue(subnetLogEntryObjectType(), elements)
}

func NewSubnetModelFromResponse(_ context.Context, sr *SubnetResponse) (*SubnetModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if sr == nil {
		diags.AddError("Constructor Error", "subnet response is nil")
		return nil, diags
	}

	model := &SubnetModel{
		Id:           types.StringValue(sr.Id),
		Name:         types.StringValue(sr.Name),
		PublicName:   types.StringValue(sr.PublicName),
		Ip:           types.StringValue(sr.Ip),
		TypeId:       types.Int64Value(sr.TypeId),
		Type:         types.StringValue(sr.Type),
		Mask:         types.Int64Value(sr.Mask),
		Gateway:      types.StringValue(sr.Gateway),
		Broadcast:    types.StringValue(sr.Broadcast),
		Dhcp:         types.Int64Value(sr.Dhcp),
		StateId:      types.Int64Value(sr.StateId),
		State:        types.StringValue(sr.State),
		DatacenterId: types.StringValue(sr.DatacenterId),
		StartDate:    types.StringValue(sr.StartDate),
	}

	if sr.Description != nil {
		model.Description = types.StringValue(*sr.Description)
	} else {
		model.Description = types.StringNull()
	}

	if sr.Subnet != nil {
		model.Subnet = types.StringValue(*sr.Subnet)
	} else {
		model.Subnet = types.StringNull()
	}

	if sr.SubnetId != nil {
		model.SubnetId = types.StringValue(*sr.SubnetId)
	} else {
		model.SubnetId = types.StringNull()
	}

	if sr.NetworkInterfaceId != nil {
		model.NetworkInterfaceId = types.StringValue(*sr.NetworkInterfaceId)
	} else {
		model.NetworkInterfaceId = types.StringNull()
	}

	if sr.ServerType != nil {
		model.ServerType = types.StringValue(*sr.ServerType)
	} else {
		model.ServerType = types.StringNull()
	}

	if sr.LoadBalancerId != nil {
		model.LoadBalancerId = types.StringValue(*sr.LoadBalancerId)
	} else {
		model.LoadBalancerId = types.StringNull()
	}

	if sr.LoadBalancerPublicName != nil {
		model.LoadBalancerPublicName = types.StringValue(*sr.LoadBalancerPublicName)
	} else {
		model.LoadBalancerPublicName = types.StringNull()
	}

	if sr.NetworkId != nil {
		model.NetworkId = types.StringValue(*sr.NetworkId)
	} else {
		model.NetworkId = types.StringNull()
	}

	if sr.NetworkPublicName != nil {
		model.NetworkPublicName = types.StringValue(*sr.NetworkPublicName)
	} else {
		model.NetworkPublicName = types.StringNull()
	}

	if sr.InverseDns != nil {
		model.InverseDns = types.StringValue(*sr.InverseDns)
	} else {
		model.InverseDns = types.StringNull()
	}

	logsList, logsDiags := NewSubnetLogsList(sr.LastLogs)
	diags.Append(logsDiags...)
	if !logsDiags.HasError() {
		model.LastLogs = logsList
	}

	return model, diags
}

func subnetLogEntryObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":     types.StringType,
			"uuid":   types.StringType,
			"date":   types.StringType,
			"action": types.StringType,
			"time":   types.Int64Type,
			"result": types.StringType,
			"type":   types.StringType,
		},
	}
}

func SubnetResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Subnet resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet IP address",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$`),
						"must be a valid IP address",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"mask": rschema.Int64Attribute{
				Required:    true,
				Description: "Subnet mask",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Required:    true,
				Description: "Datacenter ID",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter ID",
					),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Subnet description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
				},
			},
			"public_name": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet public name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"type_id": rschema.Int64Attribute{
				Computed:    true,
				Description: "Subnet type ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"type": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet": rschema.StringAttribute{
				Computed:    true,
				Description: "Subnet range",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"subnet_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Parent subnet ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_interface_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Network interface ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_type": rschema.StringAttribute{
				Computed:    true,
				Description: "Server type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Load balancer ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"load_balancer_public_name": rschema.StringAttribute{
				Computed:    true,
				Description: "Load balancer public name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": rschema.StringAttribute{
				Computed:    true,
				Description: "Network ID",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_public_name": rschema.StringAttribute{
				Computed:    true,
				Description: "Network public name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"gateway": rschema.StringAttribute{
				Computed:    true,
				Description: "Gateway IP",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"broadcast": rschema.StringAttribute{
				Computed:    true,
				Description: "Broadcast IP",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"inversedns": rschema.StringAttribute{
				Computed:    true,
				Description: "Inverse DNS",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dhcp": rschema.Int64Attribute{
				Computed:    true,
				Description: "DHCP enabled",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"state_id": rschema.Int64Attribute{
				Computed:    true,
				Description: "State ID",
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"state": rschema.StringAttribute{
				Computed:    true,
				Description: "State",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"start_date": rschema.StringAttribute{
				Computed:    true,
				Description: "Start date",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"last_logs": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "Last logs",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: map[string]rschema.Attribute{
						"id": rschema.StringAttribute{
							Computed:    true,
							Description: "Log ID",
						},
						"uuid": rschema.StringAttribute{
							Computed:    true,
							Description: "Log UUID",
						},
						"date": rschema.StringAttribute{
							Computed:    true,
							Description: "Log date",
						},
						"action": rschema.StringAttribute{
							Computed:    true,
							Description: "Log action",
						},
						"time": rschema.Int64Attribute{
							Computed:    true,
							Description: "Log time",
						},
						"result": rschema.StringAttribute{
							Computed:    true,
							Description: "Log result",
						},
						"type": rschema.StringAttribute{
							Computed:    true,
							Description: "Log type",
						},
					},
				},
			},
		},
	}
}

func SubnetDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for subnet information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Subnet identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid subnet ID",
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
				Description: "Id of the subnet to which the subnet belongs",
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
				Description: "Current state of the subnet (ACTIVE, etc.)",
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
