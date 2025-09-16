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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PublicNetworkModel struct {
	Id                types.String `tfsdk:"id"`
	PublicName        types.String `tfsdk:"public_name"`
	Description       types.String `tfsdk:"description"`
	DatacenterId      types.String `tfsdk:"datacenter_id"`
	StartDate         types.Int64  `tfsdk:"start_date"`
	SameVlan          types.Bool   `tfsdk:"same_vlan"`
	Type              types.String `tfsdk:"type"`
	State             types.String `tfsdk:"state"`
	Servers           types.List   `tfsdk:"servers"`
	Ips               types.List   `tfsdk:"ips"`
	AvailabilityZones types.List   `tfsdk:"availability_zones"`
	LastLogs          types.List   `tfsdk:"last_logs"`
}

type PublicNetworkResponse struct {
	Id                string                        `json:"id"`
	PublicName        string                        `json:"public_name"`
	Description       *string                       `json:"description"`
	DatacenterId      string                        `json:"datacenter_id"`
	StartDate         int64                         `json:"start_date"`
	SameVlan          bool                          `json:"same_vlan"`
	Type              string                        `json:"type"`
	State             string                        `json:"state"`
	Servers           []PublicNetworkServerResponse `json:"servers"`
	Ips               []string                      `json:"ips"`
	AvailabilityZones []AvailabilityZoneResponse    `json:"availability_zones"`
	LastLogs          []PublicNetworkLogResponse    `json:"last_logs"`
}

type PublicNetworkServerResponse struct {
	Id     string `json:"id"`
	VlanId int    `json:"vlan_id"`
	Mac    string `json:"mac"`
	Tagged bool   `json:"tagged"`
}

type AvailabilityZoneResponse struct {
	Id     string `json:"id"`
	VlanId int    `json:"vlan_id"`
}

type PublicNetworkLogResponse struct {
	Id     string `json:"id"`
	UUID   string `json:"uuid"`
	Date   string `json:"date"`
	Action string `json:"action"`
	Time   int    `json:"time"`
	Result string `json:"result"`
	Type   string `json:"type"`
}

func publicNetworkServerAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":      types.StringType,
		"vlan_id": types.Int64Type,
		"mac":     types.StringType,
		"tagged":  types.BoolType,
	}
}

func publicNetworkServerObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: publicNetworkServerAttributeTypes(),
	}
}

func availabilityZoneAttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"id":      types.StringType,
		"vlan_id": types.Int64Type,
	}
}

func availabilityZoneObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: availabilityZoneAttributeTypes(),
	}
}

func publicNetworkLogObjectType() types.ObjectType {
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

func NewPublicNetworkServerObject(pn PublicNetworkServerResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		publicNetworkServerAttributeTypes(),
		map[string]attr.Value{
			"id":      types.StringValue(pn.Id),
			"vlan_id": types.Int64Value(int64(pn.VlanId)),
			"mac":     types.StringValue(pn.Mac),
			"tagged":  types.BoolValue(pn.Tagged),
		},
	)
}

func NewAvailabilityZoneObject(az AvailabilityZoneResponse) (types.Object, diag.Diagnostics) {
	return types.ObjectValue(
		availabilityZoneAttributeTypes(),
		map[string]attr.Value{
			"id":      types.StringValue(az.Id),
			"vlan_id": types.Int64Value(int64(az.VlanId)),
		},
	)
}

func NewPublicNetworkServerList(pn []PublicNetworkServerResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(pn) == 0 {
		return types.ListValue(publicNetworkServerObjectType(), []attr.Value{})
	}

	values := make([]attr.Value, 0, len(pn))
	for _, server := range pn {
		serverObj, serverDiags := NewPublicNetworkServerObject(server)
		diags.Append(serverDiags...)
		if !serverDiags.HasError() {
			values = append(values, serverObj)
		}
	}

	return types.ListValue(publicNetworkServerObjectType(), values)
}

func NewAvailabilityZoneList(azs []AvailabilityZoneResponse) (types.List, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if len(azs) == 0 {
		return types.ListValue(availabilityZoneObjectType(), []attr.Value{})
	}

	values := make([]attr.Value, 0, len(azs))
	for _, az := range azs {
		azObj, azDiags := NewAvailabilityZoneObject(az)
		diags.Append(azDiags...)
		if !azDiags.HasError() {
			values = append(values, azObj)
		}
	}

	return types.ListValue(availabilityZoneObjectType(), values)
}

func NewPublicNetworkModel(ctx context.Context, pn *PublicNetworkResponse) (*PublicNetworkModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if pn == nil {
		diags.AddError("Constructor Error", "public network response is nil")
		return nil, diags
	}

	model := &PublicNetworkModel{}

	model.Id = types.StringValue(pn.Id)
	model.PublicName = types.StringValue(pn.PublicName)
	model.DatacenterId = types.StringValue(pn.DatacenterId)
	model.StartDate = types.Int64Value(pn.StartDate)
	model.SameVlan = types.BoolValue(pn.SameVlan)
	model.Type = types.StringValue(pn.Type)
	model.State = types.StringValue(pn.State)

	if pn.Description != nil {
		model.Description = types.StringValue(*pn.Description)
	} else {
		model.Description = types.StringNull()
	}

	serversList, serversDiags := NewPublicNetworkServerList(pn.Servers)
	diags.Append(serversDiags...)
	if !serversDiags.HasError() {
		model.Servers = serversList
	}

	if len(pn.Ips) > 0 {
		ipValues := make([]attr.Value, 0, len(pn.Ips))
		for _, ip := range pn.Ips {
			ipValues = append(ipValues, types.StringValue(ip))
		}
		ipsList, ipsDiags := types.ListValue(types.StringType, ipValues)
		diags.Append(ipsDiags...)
		if !ipsDiags.HasError() {
			model.Ips = ipsList
		}
	} else {
		emptyIpsList, emptyIpsDiags := types.ListValue(types.StringType, []attr.Value{})
		diags.Append(emptyIpsDiags...)
		if !emptyIpsDiags.HasError() {
			model.Ips = emptyIpsList
		}
	}

	availabilityZonesList, azDiags := NewAvailabilityZoneList(pn.AvailabilityZones)
	diags.Append(azDiags...)
	if !azDiags.HasError() {
		model.AvailabilityZones = availabilityZonesList
	}

	if len(pn.LastLogs) > 0 {
		logValues := make([]attr.Value, 0, len(pn.LastLogs))
		for _, logEntry := range pn.LastLogs {
			logObj, err := types.ObjectValue(
				map[string]attr.Type{
					"id":     types.StringType,
					"uuid":   types.StringType,
					"date":   types.StringType,
					"action": types.StringType,
					"time":   types.Int64Type,
					"result": types.StringType,
					"type":   types.StringType,
				},
				map[string]attr.Value{
					"id":     types.StringValue(logEntry.Id),
					"uuid":   types.StringValue(logEntry.UUID),
					"date":   types.StringValue(logEntry.Date),
					"action": types.StringValue(logEntry.Action),
					"time":   types.Int64Value(int64(logEntry.Time)),
					"result": types.StringValue(logEntry.Result),
					"type":   types.StringValue(logEntry.Type),
				},
			)
			if err.HasError() {
				diags.Append(err...)
				continue
			}
			logValues = append(logValues, logObj)
		}

		logsList, logsDiags := types.ListValue(publicNetworkLogObjectType(), logValues)
		diags.Append(logsDiags...)
		if !logsDiags.HasError() {
			model.LastLogs = logsList
		}
	} else {
		emptyList, emptyDiags := types.ListValueFrom(ctx, publicNetworkLogObjectType(), []PublicNetworkLogResponse{})
		diags.Append(emptyDiags...)
		if !emptyDiags.HasError() {
			model.LastLogs = emptyList
		}
	}

	return model, diags
}

func publicNetworkObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                 types.StringType,
			"public_name":        types.StringType,
			"description":        types.StringType,
			"datacenter_id":      types.StringType,
			"start_date":         types.Int64Type,
			"same_vlan":          types.BoolType,
			"type":               types.StringType,
			"state":              types.StringType,
			"servers":            types.ListType{ElemType: publicNetworkServerObjectType()},
			"ips":                types.ListType{ElemType: types.StringType},
			"availability_zones": types.ListType{ElemType: availabilityZoneObjectType()},
			"last_logs":          types.ListType{ElemType: publicNetworkLogObjectType()},
		},
	}
}

func publicNetworkNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := PublicNetworkDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Public network identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func NewPublicNetworkFromList(ctx context.Context, pnList []PublicNetworkResponse) ([]PublicNetworkModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []PublicNetworkModel

	if len(pnList) == 0 {
		return []PublicNetworkModel{}, diags
	}

	for i, pn := range pnList {
		model, modelDiags := NewPublicNetworkModel(ctx, &pn)
		if modelDiags.HasError() {
			diags.AddError(
				"List Constructor Error",
				fmt.Sprintf("Failed to create model for item %d: %s", i, modelDiags.Errors()[0].Summary()),
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

func PublicNetworkDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for public network information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Public network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid public network ID",
					),
				},
			},
			"public_name": schema.StringAttribute{
				Computed:    true,
				Description: "Public network name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Public network description",
			},
			"datacenter_id": schema.StringAttribute{
				Computed:    true,
				Description: "Datacenter identifier",
			},
			"start_date": schema.Int64Attribute{
				Computed:    true,
				Description: "Start date timestamp in milliseconds",
			},
			"same_vlan": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether the network uses the same VLAN",
			},
			"type": schema.StringAttribute{
				Computed:    true,
				Description: "Network type",
			},
			"state": schema.StringAttribute{
				Computed:    true,
				Description: "Public network state",
			},
			"servers": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of servers in the public network",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Server identifier",
						},
						"vlan_id": schema.Int64Attribute{
							Computed:    true,
							Description: "VLAN identifier",
						},
						"mac": schema.StringAttribute{
							Computed:    true,
							Description: "MAC address",
						},
						"tagged": schema.BoolAttribute{
							Computed:    true,
							Description: "Whether the server is tagged",
						},
					},
				},
			},
			"ips": schema.ListAttribute{
				Computed:    true,
				Description: "List of IP identifiers in the public network",
				ElementType: types.StringType,
			},
			"availability_zones": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of availability zones",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Availability zone identifier",
						},
						"vlan_id": schema.Int64Attribute{
							Computed:    true,
							Description: "VLAN identifier",
						},
					},
				},
			},
			"last_logs": schema.ListNestedAttribute{
				Computed:    true,
				Description: "List of recent log events",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed:    true,
							Description: "Log entry identifier",
						},
						"uuid": schema.StringAttribute{
							Computed:    true,
							Description: "Log entry UUID",
						},
						"date": schema.StringAttribute{
							Computed:    true,
							Description: "Log entry date",
						},
						"action": schema.StringAttribute{
							Computed:    true,
							Description: "Action performed",
						},
						"time": schema.Int64Attribute{
							Computed:    true,
							Description: "Time taken for the action",
						},
						"result": schema.StringAttribute{
							Computed:    true,
							Description: "Result of the action",
						},
						"type": schema.StringAttribute{
							Computed:    true,
							Description: "Log entry type",
						},
					},
				},
			},
		},
	}
}

type PublicNetworkCreateRequest struct {
	PublicName   string `json:"public_name"`
	Description  string `json:"description"`
	DatacenterId string `json:"datacenter_id"`
}

func (m *PublicNetworkModel) ToCreateRequest() PublicNetworkCreateRequest {
	return PublicNetworkCreateRequest{
		PublicName:   m.PublicName.ValueString(),
		Description:  m.Description.ValueString(),
		DatacenterId: m.DatacenterId.ValueString(),
	}
}

type PublicNetworkCreateResponse struct {
	Sync   bool                  `json:"sync"`
	Data   PublicNetworkResponse `json:"data"`
	TaskID string                `json:"task_id"`
}

func PublicNetworkResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Public network resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Public network identifier",
			},
			"public_name": rschema.StringAttribute{
				Required:    true,
				Description: "Public network name",
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
				Description: "Public network description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxDescriptionLength),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Datacenter identifier where the network will be created",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter_id",
					),
				},
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
			"servers": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "List of servers in the public network",
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

type PublicNetworkUpdateRequest struct {
	PublicName  string `json:"public_name"`
	Description string `json:"description"`
}

func (m *PublicNetworkModel) ToUpdateRequest() PublicNetworkUpdateRequest {
	return PublicNetworkUpdateRequest{
		PublicName:  m.PublicName.ValueString(),
		Description: m.Description.ValueString(),
	}
}

func (m *PublicNetworkModel) GetState() string {
	return m.State.ValueString()
}
