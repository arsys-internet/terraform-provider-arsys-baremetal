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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/server"
	"terraform-provider-arsys-baremetal/internal/util"
)

//func (m *ServerResourceModel) ReadPath() (string, diag.Diagnostics) {
//	var diags diag.Diagnostics
//	path := "/servers/" + m.ID.ValueString()
//	if len(path) == 0 || m.ID.IsNull() {
//		diags.AddError("No read path defined for model", fmt.Sprintf("Type:%T, Model:%v", m, m))
//		return path, diags
//	}
//	return path, diags
//}

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
	Hypervisor       types.String `tfsdk:"hypervisor"`
	Hostname         types.String `tfsdk:"hostname"`
	ConnectionSpeed  types.Object `tfsdk:"connection_speed"` // CORREGIDO: agregado
	Redundancy       types.Object `tfsdk:"redundancy"`
	RSAKey           types.Bool   `tfsdk:"rsa_key"`
	Snapshot         types.Object `tfsdk:"snapshot"`
	PrivateNetworks  types.List   `tfsdk:"private_networks"`
}

type ServerResourceModel struct {
	ServerModel

	ApplianceID            types.String  `tfsdk:"appliance_id"`
	DatacenterID           types.String  `tfsdk:"datacenter_id"`
	Password               types.String  `tfsdk:"password"`
	PowerOn                types.Bool    `tfsdk:"power_on"`
	FirewallPolicyID       types.String  `tfsdk:"firewall_policy_id"`
	IPID                   types.String  `tfsdk:"ip_id"`
	LoadBalancerID         types.String  `tfsdk:"load_balancer_id"`
	MonitoringPolicyID     types.String  `tfsdk:"monitoring_policy_id"`
	PrivateNetworkID       types.String  `tfsdk:"private_network_id"`
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
	FirstPassword    *string                                `json:"first_password"`
	Managed          bool                                   `json:"managed"`
	Status           server.StatusResponse                  `json:"status"`
	IPs              []server.ServersIPResponse             `json:"ips"`
	SSHPassword      bool                                   `json:"ssh_password"`
	Image            IdentifierResponse                     `json:"image"`
	Hardware         server.HardwareResponse                `json:"hardware"`
	DVD              *IdentifierResponse                    `json:"dvd,omitempty"`
	Alerts           *server.AlertResponse                  `json:"alerts,omitempty"`
	MonitoringPolicy *IdentifierResponse                    `json:"monitoring_policy,omitempty"`
	CloudPanelID     *string                                `json:"cloudpanel_id,omitempty"`
	ServerType       string                                 `json:"server_type"`
	Hypervisor       *string                                `json:"hypervisor,omitempty"`
	Hostname         string                                 `json:"hostname"`
	ConnectionSpeed  *server.ConnectionSpeedResponse        `json:"connection_speed,omitempty"`
	Redundancy       *server.RedundancyResponse             `json:"redundancy,omitempty"`
	RSAKey           interface{}                            `json:"rsa_key"`
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

func newServerModelFromResponse(ctx context.Context, sr *ServerResponse) (*ServerModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if sr == nil {
		diags.AddError("Constructor Error", "server response is nil")
		return nil, diags
	}

	model := &ServerModel{}

	model.ID = types.StringValue(sr.ID)
	model.Name = types.StringValue(sr.Name)
	model.CreationDate = types.StringValue(sr.CreationDate)
	model.Managed = types.BoolValue(sr.Managed)
	model.SSHPassword = types.BoolValue(sr.SSHPassword)
	model.ServerType = types.StringValue(sr.ServerType)
	model.Hostname = types.StringValue(sr.Hostname)

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

	if sr.Hypervisor != nil {
		model.Hypervisor = types.StringValue(*sr.Hypervisor)
	} else {
		model.Hypervisor = types.StringNull()
	}

	if sr.ConnectionSpeed != nil {
		connectionObj, connectionDiags := server.NewConnectionSpeedObject(*sr.ConnectionSpeed)
		diags.Append(connectionDiags...)
		if !connectionDiags.HasError() {
			model.ConnectionSpeed = connectionObj
			tflog.Info(ctx, fmt.Sprintf("🔍 BASE MODEL - ConnectionSpeed after NewConnectionSpeedObject: %v", connectionObj))
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
		snapshotObj, snapshotDiags := server.NewSnapshotObject(sr.Snapshot)
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

func NewServerResourceModelFromCreate(ctx context.Context, sr *ServerResponse, plan *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	baseModel, baseDiags := newServerModelFromResponse(ctx, sr)
	if baseDiags.HasError() {
		diags.Append(baseDiags...)
		return nil, diags
	}

	model := &ServerResourceModel{
		ServerModel: *baseModel,
	}

	model.ApplianceID = plan.ApplianceID
	model.DatacenterID = plan.DatacenterID

	hardwareFromAPI := true // Por defecto usar hardware de la API

	if !plan.Hardware.IsNull() && !plan.Hardware.IsUnknown() {
		hardwareAttrs := plan.Hardware.Attributes()

		// Verificar campos críticos del hardware
		allHardwareFieldsKnown := true

		if vcore, exists := hardwareAttrs["vcore"]; exists && vcore.IsUnknown() {
			allHardwareFieldsKnown = false
			tflog.Debug(ctx, "Hardware vcore is unknown in plan")
		}
		if ram, exists := hardwareAttrs["ram"]; exists && ram.IsUnknown() {
			allHardwareFieldsKnown = false
			tflog.Debug(ctx, "Hardware ram is unknown in plan")
		}
		if hdds, exists := hardwareAttrs["hdds"]; exists && hdds.IsUnknown() {
			allHardwareFieldsKnown = false
			tflog.Debug(ctx, "Hardware hdds is unknown in plan")
		}
		if cores, exists := hardwareAttrs["cores_per_processor"]; exists && cores.IsUnknown() {
			allHardwareFieldsKnown = false
			tflog.Debug(ctx, "Hardware cores_per_processor is unknown in plan")
		}

		if allHardwareFieldsKnown {
			model.Hardware = plan.Hardware
			hardwareFromAPI = false
			tflog.Info(ctx, "Using hardware from plan (all fields known)")
		}
	}

	if hardwareFromAPI {
		tflog.Info(ctx, "Using hardware from API (plan has unknown fields)")
	}

	if !plan.Password.IsUnknown() {
		model.Password = plan.Password
	} else {
		model.Password = types.StringNull()
	}

	// Campos booleanos con defaults
	if !plan.PowerOn.IsUnknown() {
		model.PowerOn = plan.PowerOn
	} else {
		model.PowerOn = types.BoolValue(true)
	}

	if !plan.BondingAllowed.IsUnknown() {
		model.BondingAllowed = plan.BondingAllowed
	} else {
		model.BondingAllowed = types.BoolValue(false)
	}

	if !plan.InstallBackupAgent.IsUnknown() {
		model.InstallBackupAgent = plan.InstallBackupAgent
	} else {
		model.InstallBackupAgent = types.BoolValue(false)
	}

	if !plan.DisabledPassword.IsUnknown() {
		model.DisabledPassword = plan.DisabledPassword
	} else {
		model.DisabledPassword = types.BoolValue(false)
	}

	if !plan.ReservedCPU.IsUnknown() {
		model.ReservedCPU = plan.ReservedCPU
	} else {
		model.ReservedCPU = types.BoolValue(false)
	}

	if !plan.ReservedRAM.IsUnknown() {
		model.ReservedRAM = plan.ReservedRAM
	} else {
		model.ReservedRAM = types.BoolValue(false)
	}

	// Campos de configuración opcionales
	if !plan.FirewallPolicyID.IsUnknown() {
		model.FirewallPolicyID = plan.FirewallPolicyID
	} else {
		model.FirewallPolicyID = types.StringNull()
	}

	if !plan.IPID.IsUnknown() {
		model.IPID = plan.IPID
	} else {
		model.IPID = types.StringNull()
	}

	if !plan.LoadBalancerID.IsUnknown() {
		model.LoadBalancerID = plan.LoadBalancerID
	} else {
		model.LoadBalancerID = types.StringNull()
	}

	if !plan.MonitoringPolicyID.IsUnknown() {
		model.MonitoringPolicyID = plan.MonitoringPolicyID
	} else {
		model.MonitoringPolicyID = types.StringNull()
	}

	if !plan.PrivateNetworkID.IsUnknown() {
		model.PrivateNetworkID = plan.PrivateNetworkID
	} else {
		model.PrivateNetworkID = types.StringNull()
	}

	if !plan.PublicKey.IsUnknown() {
		model.PublicKey = plan.PublicKey
	} else {
		model.PublicKey = types.StringNull()
	}

	if !plan.ExecutionGroup.IsUnknown() {
		model.ExecutionGroup = plan.ExecutionGroup
	} else {
		model.ExecutionGroup = types.StringNull()
	}

	if !plan.UserData.IsUnknown() {
		model.UserData = plan.UserData
	} else {
		model.UserData = types.StringNull()
	}

	if !plan.UserDataContentType.IsUnknown() {
		model.UserDataContentType = plan.UserDataContentType
	} else {
		model.UserDataContentType = types.StringNull()
	}

	if !plan.PublicConnectionSpeed.IsUnknown() {
		model.PublicConnectionSpeed = plan.PublicConnectionSpeed
	} else {
		model.PublicConnectionSpeed = types.Float64Null()
	}

	if !plan.PrivateConnectionSpeed.IsUnknown() {
		model.PrivateConnectionSpeed = plan.PrivateConnectionSpeed
	} else {
		model.PrivateConnectionSpeed = types.Float64Null()
	}

	if !plan.BackupAgentTenantData.IsUnknown() {
		model.BackupAgentTenantData = plan.BackupAgentTenantData
	} else {
		model.BackupAgentTenantData = types.StringNull()
	}

	if !plan.BackupSize.IsUnknown() {
		model.BackupSize = plan.BackupSize
	} else {
		model.BackupSize = types.StringNull()
	}

	if !plan.AvailabilityZoneID.IsUnknown() {
		model.AvailabilityZoneID = plan.AvailabilityZoneID
	} else {
		model.AvailabilityZoneID = types.StringNull()
	}

	if !plan.BackupPolicyID.IsUnknown() {
		model.BackupPolicyID = plan.BackupPolicyID
	} else {
		model.BackupPolicyID = types.StringNull()
	}

	if !plan.ConfigurationID.IsUnknown() {
		model.ConfigurationID = plan.ConfigurationID
	} else {
		model.ConfigurationID = types.StringNull()
	}

	if !plan.MonitorID.IsUnknown() {
		model.MonitorID = plan.MonitorID
	} else {
		model.MonitorID = types.StringNull()
	}

	if !plan.FirewallID.IsUnknown() {
		model.FirewallID = plan.FirewallID
	} else {
		model.FirewallID = types.StringNull()
	}

	if !plan.SSHIDs.IsUnknown() {
		model.SSHIDs = plan.SSHIDs
	} else {
		model.SSHIDs = types.ListNull(types.StringType)
	}

	if !plan.SSHPublicKeys.IsUnknown() {
		model.SSHPublicKeys = plan.SSHPublicKeys
	} else {
		model.SSHPublicKeys = types.ListNull(types.StringType)
	}

	if !plan.Locale.IsUnknown() {
		model.Locale = plan.Locale
	} else {
		model.Locale = types.StringNull()
	}

	return model, diags
}
func NewServerResourceModelFromRead(ctx context.Context, sr *ServerResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if currentState == nil {
		baseModel, baseDiags := newServerModelFromResponse(ctx, sr)
		if baseDiags.HasError() {
			diags.Append(baseDiags...)
			return nil, diags
		}

		model := &ServerResourceModel{
			ServerModel: *baseModel,
		}

		model.Password = types.StringNull()
		model.PowerOn = types.BoolValue(true)
		model.FirewallPolicyID = types.StringNull()
		model.IPID = types.StringNull()
		model.LoadBalancerID = types.StringNull()
		model.MonitoringPolicyID = types.StringNull()
		model.PrivateNetworkID = types.StringNull()
		model.PublicKey = types.StringNull()
		model.ExecutionGroup = types.StringNull()
		model.UserData = types.StringNull()
		model.UserDataContentType = types.StringNull()
		model.PublicConnectionSpeed = types.Float64Null()
		model.PrivateConnectionSpeed = types.Float64Null()
		model.BondingAllowed = types.BoolValue(false)
		model.InstallBackupAgent = types.BoolValue(false)
		model.BackupAgentTenantData = types.StringNull()
		model.BackupSize = types.StringNull()
		model.AvailabilityZoneID = types.StringNull()
		model.DisabledPassword = types.BoolValue(false)
		model.ReservedCPU = types.BoolValue(false)
		model.ReservedRAM = types.BoolValue(false)
		model.BackupPolicyID = types.StringNull()
		model.ConfigurationID = types.StringNull()
		model.MonitorID = types.StringNull()
		model.FirewallID = types.StringNull()
		model.SSHIDs = types.ListNull(types.StringType)
		model.SSHPublicKeys = types.ListNull(types.StringType)
		model.Locale = types.StringNull()

		return model, diags
	}

	model := *currentState

	model.Name = types.StringValue(sr.Name)
	if sr.Description != nil {
		model.Description = types.StringValue(*sr.Description)
	} else {
		model.Description = types.StringNull()
	}

	currentStatus := currentState.GetState()
	if sr.Status.State != currentStatus {
		statusObj, statusDiags := server.NewStatusObject(sr.Status)
		diags.Append(statusDiags...)
		if !statusDiags.HasError() {
			model.Status = statusObj
		}
	}

	if sr.FirstPassword != nil {
		newPassword := *sr.FirstPassword
		if currentState.FirstPassword.IsNull() || newPassword != currentState.FirstPassword.ValueString() {
			model.FirstPassword = types.StringValue(newPassword)
		}
	}

	if sr.Hostname != currentState.Hostname.ValueString() {
		model.Hostname = types.StringValue(sr.Hostname)
	}

	if sr.CloudPanelID != nil {
		model.CloudPanelID = types.StringValue(*sr.CloudPanelID)
	} else {
		model.CloudPanelID = types.StringNull()
	}

	if sr.Hypervisor != nil {
		model.Hypervisor = types.StringValue(*sr.Hypervisor)
	} else {
		model.Hypervisor = types.StringNull()
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

	if sr.DVD != nil {
		dvdObj, dvdDiags := NewIdentifierObject(*sr.DVD)
		diags.Append(dvdDiags...)
		if !dvdDiags.HasError() {
			model.DVD = dvdObj
		}
	} else {
		model.DVD = types.ObjectNull(IdentifierObjectType().AttrTypes)
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

	if sr.Snapshot != nil {
		snapshotObj, snapshotDiags := server.NewSnapshotObject(sr.Snapshot)
		diags.Append(snapshotDiags...)
		if !snapshotDiags.HasError() {
			model.Snapshot = snapshotObj
		}
	} else {
		model.Snapshot = types.ObjectNull(server.SnapshotObjectType().AttrTypes)
	}

	return &model, diags
}

func NewServerResourceModelFromUpdate(_ context.Context, sr *ServerResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := *currentState

	model.Name = types.StringValue(sr.Name)
	if sr.Description != nil {
		model.Description = types.StringValue(*sr.Description)
	} else {
		model.Description = types.StringNull()
	}

	return &model, diags
}

//func NewServerResourceModel(ctx context.Context, sr *ServerResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
//	baseModel, diags := newServerModelFromResponse(ctx, sr)
//	if diags.HasError() {
//		return nil, diags
//	}
//
//	resourceModel := &ServerResourceModel{
//		ServerModel: *baseModel,
//	}
//
//	if currentState != nil {
//		resourceModel.ConnectionSpeed = currentState.ConnectionSpeed
//		resourceModel.ApplianceID = currentState.ApplianceID
//		resourceModel.DatacenterID = currentState.DatacenterID
//		resourceModel.Password = currentState.Password
//		resourceModel.FirewallPolicyID = currentState.FirewallPolicyID
//		resourceModel.IPID = currentState.IPID
//		resourceModel.LoadBalancerID = currentState.LoadBalancerID
//		resourceModel.MonitoringPolicyID = currentState.MonitoringPolicyID
//		resourceModel.PrivateNetworkID = currentState.PrivateNetworkID
//		resourceModel.PublicKey = currentState.PublicKey
//		resourceModel.ExecutionGroup = currentState.ExecutionGroup
//		resourceModel.UserData = currentState.UserData
//		resourceModel.UserDataContentType = currentState.UserDataContentType
//		resourceModel.PublicConnectionSpeed = currentState.PublicConnectionSpeed
//		resourceModel.PrivateConnectionSpeed = currentState.PrivateConnectionSpeed
//		resourceModel.BackupAgentTenantData = currentState.BackupAgentTenantData
//		resourceModel.BackupSize = currentState.BackupSize
//		resourceModel.AvailabilityZoneID = currentState.AvailabilityZoneID
//		resourceModel.BackupPolicyID = currentState.BackupPolicyID
//		resourceModel.ConfigurationID = currentState.ConfigurationID
//		resourceModel.MonitorID = currentState.MonitorID
//		resourceModel.FirewallID = currentState.FirewallID
//		resourceModel.SSHIDs = currentState.SSHIDs
//		resourceModel.SSHPublicKeys = currentState.SSHPublicKeys
//		resourceModel.Locale = currentState.Locale
//		resourceModel.Name = currentState.Name
//		resourceModel.Description = currentState.Description
//		resourceModel.Hardware = currentState.Hardware
//		resourceModel.ID = currentState.ID
//		resourceModel.CreationDate = currentState.CreationDate
//		resourceModel.FirstPassword = currentState.FirstPassword
//		resourceModel.Managed = currentState.Managed
//		resourceModel.Status = currentState.Status
//		resourceModel.IPs = currentState.IPs
//		resourceModel.SSHPassword = currentState.SSHPassword
//		resourceModel.Image = currentState.Image
//		resourceModel.DVD = currentState.DVD
//		resourceModel.Alerts = currentState.Alerts
//		resourceModel.MonitoringPolicy = currentState.MonitoringPolicy
//		resourceModel.CloudPanelID = currentState.CloudPanelID
//		resourceModel.ServerType = currentState.ServerType
//		resourceModel.Hypervisor = currentState.Hypervisor
//		resourceModel.Hostname = currentState.Hostname
//		resourceModel.ConnectionSpeed = currentState.ConnectionSpeed
//		resourceModel.Redundancy = currentState.Redundancy
//		resourceModel.RSAKey = currentState.RSAKey
//		resourceModel.Snapshot = currentState.Snapshot
//		resourceModel.PrivateNetworks = currentState.PrivateNetworks
//		resourceModel.Datacenter = currentState.Datacenter
//
//		if !currentState.PowerOn.IsUnknown() {
//			resourceModel.PowerOn = currentState.PowerOn
//		} else {
//			resourceModel.PowerOn = types.BoolValue(true)
//		}
//
//		if !currentState.BondingAllowed.IsUnknown() {
//			resourceModel.BondingAllowed = currentState.BondingAllowed
//		} else {
//			resourceModel.BondingAllowed = types.BoolValue(false)
//		}
//
//		if !currentState.InstallBackupAgent.IsUnknown() {
//			resourceModel.InstallBackupAgent = currentState.InstallBackupAgent
//		} else {
//			resourceModel.InstallBackupAgent = types.BoolValue(false)
//		}
//
//		if !currentState.DisabledPassword.IsUnknown() {
//			resourceModel.DisabledPassword = currentState.DisabledPassword
//		} else {
//			resourceModel.DisabledPassword = types.BoolValue(false)
//		}
//
//		if !currentState.ReservedCPU.IsUnknown() {
//			resourceModel.ReservedCPU = currentState.ReservedCPU
//		} else {
//			resourceModel.ReservedCPU = types.BoolValue(false)
//		}
//
//		if !currentState.ReservedRAM.IsUnknown() {
//			resourceModel.ReservedRAM = currentState.ReservedRAM
//		} else {
//			resourceModel.ReservedRAM = types.BoolValue(false)
//		}
//
//	} else {
//		resourceModel.ApplianceID = types.StringValue(sr.Image.ID)
//		resourceModel.DatacenterID = types.StringValue(sr.Datacenter.ID)
//		resourceModel.Password = types.StringNull()
//		resourceModel.PowerOn = types.BoolValue(true)
//		resourceModel.FirewallPolicyID = types.StringNull()
//		resourceModel.IPID = types.StringNull()
//		resourceModel.LoadBalancerID = types.StringNull()
//		resourceModel.MonitoringPolicyID = types.StringNull()
//		resourceModel.PrivateNetworkID = types.StringNull()
//		resourceModel.PublicKey = types.StringNull()
//		resourceModel.ExecutionGroup = types.StringNull()
//		resourceModel.UserData = types.StringNull()
//		resourceModel.UserDataContentType = types.StringNull()
//		resourceModel.PublicConnectionSpeed = types.Float64Null()
//		resourceModel.PrivateConnectionSpeed = types.Float64Null()
//		resourceModel.BondingAllowed = types.BoolValue(false)
//		resourceModel.InstallBackupAgent = types.BoolValue(false)
//		resourceModel.BackupAgentTenantData = types.StringNull()
//		resourceModel.BackupSize = types.StringNull()
//		resourceModel.AvailabilityZoneID = types.StringNull()
//		resourceModel.DisabledPassword = types.BoolValue(false)
//		resourceModel.ReservedCPU = types.BoolValue(false)
//		resourceModel.ReservedRAM = types.BoolValue(false)
//		resourceModel.BackupPolicyID = types.StringNull()
//		resourceModel.ConfigurationID = types.StringNull()
//		resourceModel.MonitorID = types.StringNull()
//		resourceModel.FirewallID = types.StringNull()
//		resourceModel.SSHIDs = types.ListNull(types.StringType)
//		resourceModel.SSHPublicKeys = types.ListNull(types.StringType)
//		resourceModel.Locale = types.StringNull()
//	}
//
//	return resourceModel, diags
//}

func (s *ServerResourceModel) ToCreateRequest() ServerCreateRequest {
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
		ServerType:   "baremetal",
		ApplianceID:  s.ApplianceID.ValueString(),
		DatacenterID: s.DatacenterID.ValueString(),
		Hardware: server.HardwareCreateRequest{
			BaremetalModelID: hardwareAttrs["baremetal_model_id"].(types.String).ValueString(),
			HDDs:             hdds,
		},
	}

	if !s.Description.IsNull() {
		req.Description = s.Description.ValueString()
	}

	if !s.SSHPassword.IsNull() {
		req.SSHPassword = s.SSHPassword.ValueBool()
	}

	if !s.PowerOn.IsNull() {
		req.PowerOn = s.PowerOn.ValueBool()
	} else {
		req.PowerOn = true
	}

	if !s.RSAKey.IsNull() {
		req.RSAKey = s.RSAKey.ValueBool()
	}

	if vcore, ok := hardwareAttrs["vcore"].(types.Int64); ok && !vcore.IsNull() {
		req.Hardware.VCore = int(vcore.ValueInt64())
	}

	if cores, ok := hardwareAttrs["cores_per_processor"].(types.Int64); ok && !cores.IsNull() {
		req.Hardware.CoresPerProcessor = int(cores.ValueInt64())
	}

	if ram, ok := hardwareAttrs["ram"].(types.Int64); ok && !ram.IsNull() {
		req.Hardware.RAM = int(ram.ValueInt64())
	}

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
		model, modelDiags := NewServerModel(ctx, &serverModel)
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
			"hypervisor":        types.StringType,
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
			"hypervisor": schema.StringAttribute{
				Computed:    true,
				Description: "Server hypervisor type",
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
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
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
				Description: "Server password",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
					stringvalidator.LengthAtMost(64),
				},
			},
			"power_on": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(true),
				Description: "Whether to power on the server after creation",
			},
			"ssh_password": rschema.BoolAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Whether SSH password authentication is enabled",
			},
			"firewall_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Firewall policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid firewall_policy_id",
					),
				},
			},
			"ip_id": rschema.StringAttribute{
				Optional:    true,
				Description: "IP identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ip_id",
					),
				},
			},
			"load_balancer_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Load balancer identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid load_balancer_id",
					),
				},
			},
			"monitoring_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Monitoring policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid monitoring_policy_id",
					),
				},
			},
			"private_network_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Private network identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid private_network_id",
					),
				},
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
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false), Description: "Whether to reserve RAM",
			},
			"backup_policy_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Backup policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid backup_policy_id",
					),
				},
			},
			"configuration_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Configuration identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid backup_policy_id",
					),
				},
			},
			"monitor_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Monitor identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid backup_policy_id",
					),
				},
			},
			"firewall_id": rschema.StringAttribute{
				Optional:    true,
				Description: "Firewall identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid backup_policy_id",
					),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_password": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"managed": rschema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"status": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: server.StatusResourceSchema(),
			},
			"ips": rschema.ListNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersIPResourceSchema(),
				},
			},
			"image": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: BaseIdentifierResourceAttributes(),
			},
			"hardware": rschema.SingleNestedAttribute{
				Required: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: server.HardwareResourceSchema(),
			},
			"dvd": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: BaseIdentifierResourceAttributes(),
			},
			"alerts": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Description: "Server alerts",
				Attributes:  server.AlertsResourceSchema(),
			},
			"monitoring_policy": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: BaseIdentifierResourceAttributes(),
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"server_type": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hypervisor": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Description: "Server hypervisor type",
			},
			"hostname": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_speed": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Description: "Connection speed configuration",
				Attributes:  server.ConnectionSpeedResourceSchema(),
			},
			"redundancy": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: server.RedundancyResourceSchema(),
			},
			"rsa_key": rschema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshot": rschema.SingleNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: server.SnapshotResourceSchema(),
			},
			"private_networks": rschema.ListNestedAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersPrivateNetworksResourceSchema(),
				},
			},
		},
	}
}
