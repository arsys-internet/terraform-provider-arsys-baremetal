package models

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	rschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/server"
	"terraform-provider-arsys-baremetal/internal/util"
)

type ServerModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Datacenter       types.Object `tfsdk:"datacenter"`
	CreationDate     types.String `tfsdk:"creation_date"`
	FirstPassword    types.String `tfsdk:"first_password"`
	Managed          types.Bool   `tfsdk:"managed"`
	Status           types.Object `tfsdk:"status"`
	IPs              types.List   `tfsdk:"ips"`
	SSHPassword      types.Bool   `tfsdk:"ssh_password"`
	Image            types.Object `tfsdk:"image"`
	Hardware         types.Object `tfsdk:"hardware"`
	DVD              types.Object `tfsdk:"dvd"`
	Alerts           types.Object `tfsdk:"alerts"`
	MonitoringPolicy types.Object `tfsdk:"monitoring_policy"`
	CloudPanelID     types.String `tfsdk:"cloudpanel_id"`
	ServerType       types.String `tfsdk:"server_type"`
	Hostname         types.String `tfsdk:"hostname"`
	ConnectionSpeed  types.Object `tfsdk:"connection_speed"` // CORREGIDO: agregado
	Redundancy       types.Object `tfsdk:"redundancy"`
	RSAKey           types.Bool   `tfsdk:"rsa_key"`
	Snapshot         types.Object `tfsdk:"snapshot"`
	PrivateNetworks  types.List   `tfsdk:"private_networks"`
}

type ServerResourceModel struct {
	ServerModel

	// Campos específicos de resource (inputs del usuario)
	ApplianceID            types.String  `tfsdk:"appliance_id"`
	DatacenterID           types.String  `tfsdk:"datacenter_id"`
	Password               types.String  `tfsdk:"password"`
	PowerOn                types.Bool    `tfsdk:"power_on"`
	FirewallPolicyID       types.String  `tfsdk:"firewall_policy_id"`
	IPID                   types.String  `tfsdk:"ip_id"`
	LoadBalancerID         types.String  `tfsdk:"load_balancer_id"`
	MonitoringPolicyID     types.String  `tfsdk:"monitoring_policy_id"`
	PrivateNetworkID       types.String  `tfsdk:"private_network_id"`
	SiteID                 types.String  `tfsdk:"site_id"`
	PublicKey              types.String  `tfsdk:"public_key"`
	ExecutionGroup         types.String  `tfsdk:"execution_group"`
	UserData               types.String  `tfsdk:"user_data"`
	UserDataContentType    types.String  `tfsdk:"user_data_content_type"`
	PublicConnectionSpeed  types.Float64 `tfsdk:"public_connection_speed"`
	PrivateConnectionSpeed types.Float64 `tfsdk:"private_connection_speed"`
	BondingAllowed         types.Bool    `tfsdk:"bonding_allowed"`
	InstallBackupAgent     types.Bool    `tfsdk:"install_backup_agent"`
	BackupAgentTenantData  types.String  `tfsdk:"backup_agent_tenant_data"`
	BackupSize             types.String  `tfsdk:"backup_size"`
	AvailabilityZoneID     types.String  `tfsdk:"availability_zone_id"`

	DisabledPassword types.Bool   `tfsdk:"disabled_password"`
	ReservedCPU      types.Bool   `tfsdk:"reserved_cpu"`
	ReservedRAM      types.Bool   `tfsdk:"reserved_ram"`
	BackupPolicyID   types.String `tfsdk:"backup_policy_id"`
	ConfigurationID  types.String `tfsdk:"configuration_id"`
	MonitorID        types.String `tfsdk:"monitor_id"`
	FirewallID       types.String `tfsdk:"firewall_id"`
	SSHIDs           types.List   `tfsdk:"ssh_ids"`
	SSHPublicKeys    types.List   `tfsdk:"ssh_public_keys"`
	Locale           types.String `tfsdk:"locale"`
}

type ServerResponse struct {
	ID               string                                 `json:"id"`
	Name             string                                 `json:"name"`
	Description      *string                                `json:"description,omitempty"`
	Datacenter       BaseDatacenterResponse                 `json:"datacenter"`
	CreationDate     string                                 `json:"creation_date"`
	FirstPassword    *string                                `json:"first_password"` // CORREGIDO: nullable
	Managed          bool                                   `json:"managed"`
	Status           server.StatusResponse                  `json:"status"`
	IPs              []server.ServersIPResponse             `json:"ips"`
	SSHPassword      bool                                   `json:"ssh_password"`
	Image            IdentifierResponse                     `json:"image"`
	Hardware         server.HardwareResponse                `json:"hardware"`
	DVD              *IdentifierResponse                    `json:"dvd,omitempty"`
	Alerts           *server.AlertResponse                  `json:"alerts,omitempty"`
	MonitoringPolicy *IdentifierResponse                    `json:"monitoring_policy,omitempty"` // CORREGIDO: nullable
	CloudPanelID     *string                                `json:"cloudpanel_id,omitempty"`     // CORREGIDO: nullable
	ServerType       string                                 `json:"server_type"`
	Hostname         string                                 `json:"hostname"`
	ConnectionSpeed  *server.ConnectionSpeedResponse        `json:"connection_speed,omitempty"` // AGREGADO
	Redundancy       *server.RedundancyResponse             `json:"redundancy,omitempty"`       // CORREGIDO: nullable
	RSAKey           interface{}                            `json:"rsa_key"`                    // Puede ser int o bool
	Snapshot         *server.SnapshotResponse               `json:"snapshot,omitempty"`
	PrivateNetworks  []server.ServersPrivateNetworkResponse `json:"private_networks"`
}

type ServerCreateRequest struct {
	Name                   string                       `json:"name"`
	Description            string                       `json:"description,omitempty"`
	ServerType             string                       `json:"server_type"`
	ApplianceID            string                       `json:"appliance_id"`
	DatacenterID           string                       `json:"datacenter_id"`
	SSHPassword            bool                         `json:"ssh_password"`
	PowerOn                bool                         `json:"power_on"`
	RSAKey                 bool                         `json:"rsa_key"`
	Hardware               server.HardwareCreateRequest `json:"hardware"`
	Password               string                       `json:"password,omitempty"`
	FirewallPolicyID       string                       `json:"firewall_policy_id,omitempty"`
	IPID                   string                       `json:"ip_id,omitempty"`
	LoadBalancerID         string                       `json:"load_balancer_id,omitempty"`
	MonitoringPolicyID     string                       `json:"monitoring_policy_id,omitempty"`
	PrivateNetworkID       string                       `json:"private_network_id,omitempty"`
	SiteID                 string                       `json:"site_id,omitempty"`
	PublicKey              string                       `json:"public_key,omitempty"`
	ExecutionGroup         string                       `json:"execution_group,omitempty"`
	UserData               string                       `json:"user_data,omitempty"`
	UserDataContentType    string                       `json:"user_data_content_type,omitempty"`
	PublicConnectionSpeed  float64                      `json:"public_connection_speed,omitempty"`
	PrivateConnectionSpeed float64                      `json:"private_connection_speed,omitempty"`
	BondingAllowed         bool                         `json:"bonding_allowed"`
	InstallBackupAgent     bool                         `json:"install_backup_agent"`
	BackupAgentTenantData  string                       `json:"backup_agent_tenant_data,omitempty"`
	BackupSize             string                       `json:"backup_size,omitempty"`
	AvailabilityZoneID     string                       `json:"availability_zone_id,omitempty"`
}

type ServerUpdateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (s *ServerResourceModel) GetState() string {
	if s == nil || s.Status.IsNull() {
		return ""
	}

	var status server.StatusModel
	diags := s.Status.As(context.Background(), &status, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return ""
	}
	return status.State.ValueString()
}

func newServerModelFromResponse(_ context.Context, sr *ServerResponse) (*ServerModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if sr == nil {
		diags.AddError("Constructor Error", "server response is nil")
		return nil, diags
	}

	model := &ServerModel{}

	// Campos básicos siempre presentes
	model.ID = types.StringValue(sr.ID)
	model.Name = types.StringValue(sr.Name)
	model.CreationDate = types.StringValue(sr.CreationDate)
	model.Managed = types.BoolValue(sr.Managed)
	model.SSHPassword = types.BoolValue(sr.SSHPassword)
	model.ServerType = types.StringValue(sr.ServerType)
	model.Hostname = types.StringValue(sr.Hostname)

	// Campos nullable
	if sr.Description != nil {
		model.Description = types.StringValue(*sr.Description)
	} else {
		model.Description = types.StringNull()
	}

	if sr.FirstPassword != nil {
		model.FirstPassword = types.StringValue(*sr.FirstPassword)
	} else {
		model.FirstPassword = types.StringNull()
	}

	if sr.CloudPanelID != nil {
		model.CloudPanelID = types.StringValue(*sr.CloudPanelID)
	} else {
		model.CloudPanelID = types.StringNull()
	}

	switch v := sr.RSAKey.(type) {
	case bool:
		model.RSAKey = types.BoolValue(v)
	case int:
		model.RSAKey = types.BoolValue(v != 0)
	case float64:
		model.RSAKey = types.BoolValue(v != 0)
	default:
		model.RSAKey = types.BoolValue(false)
	}

	// Objetos complejos siempre presentes
	datacenterObj, datacenterDiags := NewBaseDatacenterObject(sr.Datacenter)
	diags.Append(datacenterDiags...)
	if !datacenterDiags.HasError() {
		model.Datacenter = datacenterObj
	}

	statusObj, statusDiags := server.NewStatusObject(sr.Status)
	diags.Append(statusDiags...)
	if !statusDiags.HasError() {
		model.Status = statusObj
	}

	imageObj, imageDiags := NewIdentifierObject(sr.Image)
	diags.Append(imageDiags...)
	if !imageDiags.HasError() {
		model.Image = imageObj
	}

	hardwareObj, hardwareDiags := server.NewHardwareObject(sr.Hardware)
	diags.Append(hardwareDiags...)
	if !hardwareDiags.HasError() {
		model.Hardware = hardwareObj
	}

	ipsList, ipsDiags := server.NewServersIPList(sr.IPs)
	diags.Append(ipsDiags...)
	if !ipsDiags.HasError() {
		model.IPs = ipsList
	}

	pnList, pnDiags := server.NewServersPrivateNetworkList(sr.PrivateNetworks)
	diags.Append(pnDiags...)
	if !pnDiags.HasError() {
		model.PrivateNetworks = pnList
	}

	// Objetos nullable
	if sr.DVD != nil {
		dvdObj, dvdDiags := NewIdentifierObject(*sr.DVD)
		diags.Append(dvdDiags...)
		if !dvdDiags.HasError() {
			model.DVD = dvdObj
		}
	} else {
		model.DVD = types.ObjectNull(IdentifierObjectType().AttrTypes)
	}

	if sr.Alerts != nil {
		alertsObj, alertsDiags := server.NewAlertsObject(sr.Alerts)
		diags.Append(alertsDiags...)
		if !alertsDiags.HasError() {
			model.Alerts = alertsObj
		}
	} else {
		model.Alerts = types.ObjectNull(server.AlertsObjectType().AttrTypes)
	}

	if sr.MonitoringPolicy != nil {
		monitoringObj, monitoringDiags := NewIdentifierObject(*sr.MonitoringPolicy)
		diags.Append(monitoringDiags...)
		if !monitoringDiags.HasError() {
			model.MonitoringPolicy = monitoringObj
		}
	} else {
		model.MonitoringPolicy = types.ObjectNull(IdentifierObjectType().AttrTypes)
	}

	if sr.ConnectionSpeed != nil {
		connectionObj, connectionDiags := server.NewConnectionSpeedObject(*sr.ConnectionSpeed)
		diags.Append(connectionDiags...)
		if !connectionDiags.HasError() {
			model.ConnectionSpeed = connectionObj
		}
	} else {
		model.ConnectionSpeed = types.ObjectNull(server.ConnectionSpeedObjectType().AttrTypes)
	}

	if sr.Redundancy != nil {
		redundancyObj, redundancyDiags := server.NewRedundancyObject(*sr.Redundancy)
		diags.Append(redundancyDiags...)
		if !redundancyDiags.HasError() {
			model.Redundancy = redundancyObj
		}
	} else {
		model.Redundancy = types.ObjectNull(server.RedundancyObjectType().AttrTypes)
	}

	if sr.Snapshot != nil {
		snapshotObj, snapshotDiags := server.NewSnapshotObject(*sr.Snapshot)
		diags.Append(snapshotDiags...)
		if !snapshotDiags.HasError() {
			model.Snapshot = snapshotObj
		}
	} else {
		model.Snapshot = types.ObjectNull(server.SnapshotObjectType().AttrTypes)
	}

	return model, diags
}

func NewServerModel(ctx context.Context, sr *ServerResponse) (*ServerModel, diag.Diagnostics) {
	return newServerModelFromResponse(ctx, sr)
}

func NewServerResourceModel(ctx context.Context, sr *ServerResponse) (*ServerResourceModel, diag.Diagnostics) {
	baseModel, diags := newServerModelFromResponse(ctx, sr)
	if diags.HasError() {
		return nil, diags
	}

	resourceModel := &ServerResourceModel{
		ServerModel: *baseModel,
		// Los campos específicos de resource se mantienen null en lectura
		ApplianceID:            types.StringNull(),
		DatacenterID:           types.StringNull(),
		Password:               types.StringNull(),
		PowerOn:                types.BoolNull(),
		FirewallPolicyID:       types.StringNull(),
		IPID:                   types.StringNull(),
		LoadBalancerID:         types.StringNull(),
		MonitoringPolicyID:     types.StringNull(),
		PrivateNetworkID:       types.StringNull(),
		SiteID:                 types.StringNull(),
		PublicKey:              types.StringNull(),
		ExecutionGroup:         types.StringNull(),
		UserData:               types.StringNull(),
		UserDataContentType:    types.StringNull(),
		PublicConnectionSpeed:  types.Float64Null(),
		PrivateConnectionSpeed: types.Float64Null(),
		BondingAllowed:         types.BoolNull(),
		InstallBackupAgent:     types.BoolNull(),
		BackupAgentTenantData:  types.StringNull(),
		BackupSize:             types.StringNull(),
		AvailabilityZoneID:     types.StringNull(),
		DisabledPassword:       types.BoolNull(),
		ReservedCPU:            types.BoolNull(),
		ReservedRAM:            types.BoolNull(),
		BackupPolicyID:         types.StringNull(),
		ConfigurationID:        types.StringNull(),
		MonitorID:              types.StringNull(),
		FirewallID:             types.StringNull(),
		SSHIDs:                 types.ListNull(types.StringType),
		SSHPublicKeys:          types.ListNull(types.StringType),
		Locale:                 types.StringNull(),
	}

	return resourceModel, diags
}

func (s *ServerResourceModel) ToCreateRequest() ServerCreateRequest {
	// Procesar hardware
	hardwareAttrs := s.Hardware.Attributes()

	var hdds []server.HDDCreateRequest
	if hddsList, ok := hardwareAttrs["hdds"].(types.List); ok {
		for _, hddVal := range hddsList.Elements() {
			if hdd, ok := hddVal.(types.Object); ok {
				hdds = append(hdds, server.HDDCreateRequest{
					Size:   int(hdd.Attributes()["size"].(types.Int64).ValueInt64()),
					IsMain: hdd.Attributes()["is_main"].(types.Bool).ValueBool(),
				})
			}
		}
	}

	req := ServerCreateRequest{
		Name:         s.Name.ValueString(),
		ServerType:   "baremetal", // Fijo para baremetal
		ApplianceID:  s.ApplianceID.ValueString(),
		DatacenterID: s.DatacenterID.ValueString(),
		Hardware: server.HardwareCreateRequest{
			BaremetalModelID: hardwareAttrs["baremetal_model_id"].(types.String).ValueString(),
			HDDs:             hdds,
		},
	}

	// Campos opcionales con valores por defecto
	if !s.Description.IsNull() {
		req.Description = s.Description.ValueString()
	}

	if !s.SSHPassword.IsNull() {
		req.SSHPassword = s.SSHPassword.ValueBool()
	}

	if !s.PowerOn.IsNull() {
		req.PowerOn = s.PowerOn.ValueBool()
	} else {
		req.PowerOn = true // Default
	}

	if !s.RSAKey.IsNull() {
		req.RSAKey = s.RSAKey.ValueBool()
	}

	// Hardware opcionales
	if vcore, ok := hardwareAttrs["vcore"].(types.Int64); ok && !vcore.IsNull() {
		req.Hardware.VCore = int(vcore.ValueInt64())
	}

	if cores, ok := hardwareAttrs["cores_per_processor"].(types.Int64); ok && !cores.IsNull() {
		req.Hardware.CoresPerProcessor = int(cores.ValueInt64())
	}

	if ram, ok := hardwareAttrs["ram"].(types.Int64); ok && !ram.IsNull() {
		req.Hardware.RAM = int(ram.ValueInt64())
	}

	// Campos específicos del SDKv2
	if !s.Password.IsNull() {
		req.Password = s.Password.ValueString()
	}

	if !s.FirewallPolicyID.IsNull() {
		req.FirewallPolicyID = s.FirewallPolicyID.ValueString()
	}

	if !s.IPID.IsNull() {
		req.IPID = s.IPID.ValueString()
	}

	if !s.LoadBalancerID.IsNull() {
		req.LoadBalancerID = s.LoadBalancerID.ValueString()
	}

	if !s.MonitoringPolicyID.IsNull() {
		req.MonitoringPolicyID = s.MonitoringPolicyID.ValueString()
	}

	if !s.PrivateNetworkID.IsNull() {
		req.PrivateNetworkID = s.PrivateNetworkID.ValueString()
	}

	if !s.SiteID.IsNull() {
		req.SiteID = s.SiteID.ValueString()
	}

	if !s.PublicKey.IsNull() {
		req.PublicKey = s.PublicKey.ValueString()
	}

	if !s.ExecutionGroup.IsNull() {
		req.ExecutionGroup = s.ExecutionGroup.ValueString()
	}

	if !s.UserData.IsNull() {
		req.UserData = s.UserData.ValueString()
	}

	if !s.UserDataContentType.IsNull() {
		req.UserDataContentType = s.UserDataContentType.ValueString()
	}

	if !s.PublicConnectionSpeed.IsNull() {
		req.PublicConnectionSpeed = s.PublicConnectionSpeed.ValueFloat64()
	}

	if !s.PrivateConnectionSpeed.IsNull() {
		req.PrivateConnectionSpeed = s.PrivateConnectionSpeed.ValueFloat64()
	}

	if !s.BondingAllowed.IsNull() {
		req.BondingAllowed = s.BondingAllowed.ValueBool()
	}

	if !s.InstallBackupAgent.IsNull() {
		req.InstallBackupAgent = s.InstallBackupAgent.ValueBool()
	}

	if !s.BackupAgentTenantData.IsNull() {
		req.BackupAgentTenantData = s.BackupAgentTenantData.ValueString()
	}

	if !s.BackupSize.IsNull() {
		req.BackupSize = s.BackupSize.ValueString()
	}

	if !s.AvailabilityZoneID.IsNull() {
		req.AvailabilityZoneID = s.AvailabilityZoneID.ValueString()
	}

	return req
}

func (m *ServerModel) ToUpdateRequest() ServerUpdateRequest {
	return ServerUpdateRequest{
		Name:        m.Name.ValueString(),
		Description: m.Description.ValueString(),
	}
}

func NewServerFromList(ctx context.Context, serverList []ServerResponse) ([]ServerModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}
	var models []ServerModel

	if len(serverList) == 0 {
		return []ServerModel{}, diags
	}

	for i, serverModel := range serverList {
		model, modelDiags := NewServerModel(ctx, &serverModel) // <-- CORREGIDO: era NewServer
		if modelDiags.HasError() {
			diags.AddError(
				"Build error",
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

func serverModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                types.StringType,
			"name":              types.StringType,
			"description":       types.StringType,
			"datacenter":        baseDatacenterObjectType(),
			"creation_date":     types.StringType,
			"first_password":    types.StringType,
			"managed":           types.BoolType,
			"status":            server.StatusObjectType(),
			"ips":               types.ListType{ElemType: server.ServersIPObjectType()},
			"ssh_password":      types.BoolType,
			"image":             IdentifierObjectType(),
			"hardware":          server.HardwareObjectType(),
			"dvd":               IdentifierObjectType(),
			"alerts":            server.AlertsObjectType(),
			"monitoring_policy": IdentifierObjectType(),
			"cloudpanel_id":     types.StringType,
			"server_type":       types.StringType,
			"hostname":          types.StringType,
			"connection_speed":  server.ConnectionSpeedObjectType(),
			"redundancy":        server.RedundancyObjectType(),
			"rsa_key":           types.BoolType,
			"snapshot":          server.SnapshotObjectType(),
			"private_networks":  types.ListType{ElemType: server.ServersPrivateNetworkObjectType()},
		},
	}
}

func serverNestedAttributeObject() schema.NestedAttributeObject {
	existingSchema := ServerDataSourceSchema(context.Background())

	attributes := make(map[string]schema.Attribute)
	for name, attribute := range existingSchema.Attributes {
		if name == "id" {
			attributes[name] = schema.StringAttribute{
				Computed:    true,
				Description: "Server identifier",
			}
		} else {
			attributes[name] = attribute
		}
	}

	return schema.NestedAttributeObject{
		Attributes: attributes,
	}
}

func ServerDataSourceSchema(_ context.Context) schema.Schema {
	return schema.Schema{
		Description: "Data source for server information",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required:    true,
				Description: "Server identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ID (e.g., 4EFAD5836CE43ACA502FD5B99BEE44EF)",
					),
				},
			},
			"name": schema.StringAttribute{
				Computed:    true,
				Description: "Server name",
			},
			"description": schema.StringAttribute{
				Computed:    true,
				Description: "Server description",
			},
			"datacenter": BaseDatacenterNestedAttribute(),
			"creation_date": schema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
			},
			"first_password": schema.StringAttribute{
				Computed:    true,
				Description: "First password generated",
			},
			"managed": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether server is managed",
			},
			"status": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server status",
				Attributes:  server.StatusDataSourceSchema(),
			},
			"ips": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Server IP addresses",
				NestedObject: schema.NestedAttributeObject{
					Attributes: server.ServersIPDataSourceSchema(),
				},
			},
			"ssh_password": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether SSH password authentication is enabled",
			},
			"image": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierAttributes(),
			},
			"hardware": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Hardware configuration",
				Attributes:  server.HardwareDataSourceSchema(),
			},
			"dvd": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierAttributes(),
			},
			"alerts": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server alerts",
				Attributes:  server.AlertsDataSourceSchema(),
			},
			"monitoring_policy": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierAttributes(),
			},
			"cloudpanel_id": schema.StringAttribute{
				Computed: true,
			},
			"server_type": schema.StringAttribute{
				Computed: true,
			},
			"hostname": schema.StringAttribute{
				Computed: true,
			},
			"connection_speed": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Connection speed configuration",
				Attributes:  server.ConnectionSpeedDataSourceSchema(),
			},
			"redundancy": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.RedundancyDataSourceSchema(),
			},
			"rsa_key": schema.BoolAttribute{
				Computed: true,
			},
			"snapshot": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.SnapshotDataSourceSchema(),
			},
			"private_networks": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: server.ServersPrivateNetworksDataSourceSchema(),
				},
			},
		},
	}
}

func ServerResourceSchema(_ context.Context) rschema.Schema {
	return rschema.Schema{
		Description: "Server resource",
		Attributes: map[string]rschema.Attribute{
			"id": rschema.StringAttribute{
				Computed:    true,
				Description: "Server identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": rschema.StringAttribute{
				Required:    true,
				Description: "Server name",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxNameLength),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server description",
			},
			"appliance_id": rschema.StringAttribute{
				Required:    true,
				Description: "Appliance identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid appliance_id",
					),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Required:    true,
				Description: "Datacenter identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter_id",
					),
				},
			},
			"password": rschema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Server password",
			},
			"power_on": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to power on the server after creation",
			},
			"ssh_password": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether SSH password authentication is enabled",
			},
			"firewall_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Firewall policy identifier",
			},
			"ip_id": rschema.StringAttribute{
				Optional:    true,
				Description: "IP identifier",
			},
			"load_balancer_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Load balancer identifier",
			},
			"monitoring_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Monitoring policy identifier",
			},
			"private_network_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Private network identifier",
			},
			"site_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Site identifier",
			},
			"public_key": rschema.StringAttribute{
				Optional:    true,
				Description: "Public key",
			},
			"execution_group": rschema.StringAttribute{
				Optional:    true,
				Description: "Execution group",
			},
			"user_data": rschema.StringAttribute{
				Optional:    true,
				Description: "User data script",
			},
			"user_data_content_type": rschema.StringAttribute{
				Optional:    true,
				Description: "User data content type",
			},
			"public_connection_speed": rschema.Float64Attribute{
				Optional:    true,
				Description: "Public connection speed",
			},
			"private_connection_speed": rschema.Float64Attribute{
				Optional:    true,
				Description: "Private connection speed",
			},
			"bonding_allowed": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether bonding is allowed",
			},
			"install_backup_agent": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to install backup agent",
			},
			"backup_agent_tenant_data": rschema.StringAttribute{
				Optional:    true,
				Description: "Backup agent tenant data",
			},
			"backup_size": rschema.StringAttribute{
				Optional:    true,
				Description: "Backup size",
			},
			"availability_zone_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Availability zone identifier",
			},

			// Campos adicionales de tu versión original
			"disabled_password": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to disable password authentication",
			},
			"reserved_cpu": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to reserve CPU",
			},
			"reserved_ram": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to reserve RAM",
			},
			"backup_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Backup policy identifier",
			},
			"configuration_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Configuration identifier",
			},
			"monitor_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Monitor identifier",
			},
			"firewall_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Firewall identifier",
			},
			"ssh_ids": rschema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "SSH key identifiers",
			},
			"ssh_public_keys": rschema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "SSH public keys",
			},
			"locale": rschema.StringAttribute{
				Optional:    true,
				Description: "Server locale",
			},

			"datacenter": BaseDatacenterResourceNestedAttribute(),
			"creation_date": rschema.StringAttribute{
				Computed: true,
			},
			"first_password": rschema.StringAttribute{
				Computed: true,
			},
			"managed": rschema.BoolAttribute{
				Computed: true,
			},
			"status": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.StatusResourceSchema(),
			},
			"ips": rschema.ListNestedAttribute{
				Computed: true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersIPResourceSchema(),
				},
			},
			"image": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierResourceAttributes(),
			},
			"hardware": rschema.SingleNestedAttribute{
				Required:   true,
				Attributes: server.HardwareResourceSchema(),
			},
			"dvd": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierResourceAttributes(),
			},
			"alerts": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server alerts",
				Attributes:  server.AlertsResourceSchema(),
			},
			"monitoring_policy": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierResourceAttributes(),
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed: true,
			},
			"server_type": rschema.StringAttribute{
				Computed: true,
			},
			"hostname": rschema.StringAttribute{
				Computed: true,
			},
			"connection_speed": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Connection speed configuration",
				Attributes:  server.ConnectionSpeedResourceSchema(),
			},
			"redundancy": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.RedundancyResourceSchema(),
			},
			"rsa_key": rschema.BoolAttribute{
				Computed: true,
			},
			"snapshot": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.SnapshotResourceSchema(),
			},
			"private_networks": rschema.ListNestedAttribute{
				Computed: true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersPrivateNetworksResourceSchema(),
				},
			},
		},
	}
}
