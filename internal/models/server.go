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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
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
	"terraform-provider-arsys-baremetal/internal/util/helper"
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

type ServerBaseModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Description      types.String `tfsdk:"description"`
	Datacenter       types.Object `tfsdk:"datacenter"`
	CreationDate     types.String `tfsdk:"creation_date"`
	FirstPassword    types.String `tfsdk:"first_password"`
	Managed          types.Bool   `tfsdk:"managed"`
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

type ServerDetailModel struct {
	ServerBaseModel
	Status           types.Object `tfsdk:"status"`
	RecoveryMode     types.Bool   `tfsdk:"recovery_mode"`
	RecoveryImageOS  types.String `tfsdk:"recovery_image_os"`
	RecoveryUser     types.String `tfsdk:"recovery_user"`
	RecoveryPassword types.String `tfsdk:"recovery_password"`
}

type ServerResourceModel struct {
	ServerDetailModel

	// Campos de configuración requeridos
	ApplianceID            types.String  `tfsdk:"appliance_id"`
	DatacenterID           types.String  `tfsdk:"datacenter_id"`
	Password               types.String  `tfsdk:"password"`
	PowerOn                types.Bool    `tfsdk:"power_on"`
	FirewallPolicyID       types.String  `tfsdk:"firewall_policy_id"`
	IPID                   types.String  `tfsdk:"ip_id"`
	LoadBalancerID         types.String  `tfsdk:"load_balancer_id"`
	MonitoringPolicyID     types.String  `tfsdk:"monitoring_policy_id"`
	PrivateNetworkID       types.String  `tfsdk:"private_network_id"`
	UserData               types.String  `tfsdk:"user_data"`
	UserDataContentType    types.String  `tfsdk:"user_data_content_type"`
	PublicConnectionSpeed  types.Float64 `tfsdk:"public_connection_speed"`
	PrivateConnectionSpeed types.Float64 `tfsdk:"private_connection_speed"`
	BondingAllowed         types.Bool    `tfsdk:"bonding_allowed"`
	InstallBackupAgent     types.Bool    `tfsdk:"install_backup_agent"`
	BackupAgentTenantData  types.String  `tfsdk:"backup_agent_tenant_data"`
	BackupSize             types.String  `tfsdk:"backup_size"`
	AvailabilityZoneID     types.String  `tfsdk:"availability_zone_id"`
	PublicKey              types.List    `tfsdk:"public_key"`
	SiteID                 types.String  `tfsdk:"site_id"`
	SSHIDs                 types.List    `tfsdk:"ssh_ids"`
	DisabledPassword       types.Bool    `tfsdk:"disabled_password"`
}

type ServerBaseResponse struct {
	ID               string                                 `json:"id"`
	Name             string                                 `json:"name"`
	Description      *string                                `json:"description,omitempty"`
	Datacenter       BaseDatacenterResponse                 `json:"datacenter"`
	CreationDate     string                                 `json:"creation_date"`
	FirstPassword    *string                                `json:"first_password"`
	Managed          bool                                   `json:"managed"`
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

type ServerDetailResponse struct {
	ServerBaseResponse
	Status           server.StatusDetailResponse `json:"status"`
	RecoveryMode     *bool                       `json:"recovery_mode,omitempty"`
	RecoveryImageOS  *string                     `json:"recovery_image_os,omitempty"`
	RecoveryUser     *string                     `json:"recovery_user,omitempty"`
	RecoveryPassword *string                     `json:"recovery_password,omitempty"`
}

type ServerCreateRequest struct {
	// REQUIRED fields - without pointer, without omitempty
	Name         string                       `json:"name"`
	ServerType   string                       `json:"server_type"`
	ApplianceID  string                       `json:"appliance_id"`
	DatacenterID string                       `json:"datacenter_id"`
	Hardware     server.HardwareCreateRequest `json:"hardware"`

	// Fields with known DEFAULTS - no pointer, no omitempty
	SSHPassword        bool `json:"ssh_password"`
	PowerOn            bool `json:"power_on"`
	RSAKey             bool `json:"rsa_key"`
	InstallBackupAgent bool `json:"install_backup_agent"`

	// OPTIONAL fields
	Description            *string  `json:"description,omitempty"`
	Password               *string  `json:"password,omitempty"`
	FirewallPolicyID       *string  `json:"firewall_policy_id,omitempty"`
	IPID                   *string  `json:"ip_id,omitempty"`
	LoadBalancerID         *string  `json:"load_balancer_id,omitempty"`
	MonitoringPolicyID     *string  `json:"monitoring_policy_id,omitempty"`
	PrivateNetworkID       *string  `json:"private_network_id,omitempty"`
	UserData               *string  `json:"user_data,omitempty"`
	UserDataContentType    *string  `json:"user_data_content_type,omitempty"`
	PublicConnectionSpeed  *float64 `json:"public_connection_speed,omitempty"`
	PrivateConnectionSpeed *float64 `json:"private_connection_speed,omitempty"`
	BondingAllowed         *bool    `json:"bonding_allowed,omitempty"`
	SiteID                 *string  `json:"site_id,omitempty"`
	DisabledPassword       *bool    `json:"disabledPassword,omitempty"`
	PublicKey              []string `json:"public_key,omitempty"`
	SSHIDs                 []string `json:"sshIds,omitempty"`

	BackupAgentTenantData *string `json:"backup_agent_tenant_data,omitempty"`
	BackupSize            *string `json:"backup_size,omitempty"`
	AvailabilityZoneID    *string `json:"availability_zone_id,omitempty"`
}
type ServerUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (s *ServerResourceModel) GetState() string {
	if s == nil || s.Status.IsNull() {
		return ""
	}

	var status server.StatusDetailModel
	diags := s.Status.As(context.Background(), &status, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return ""
	}
	return status.State.ValueString()
}

func newServerBaseModelFromResponse(_ context.Context, sr *ServerBaseResponse) (*ServerBaseModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	if sr == nil {
		diags.AddError("Constructor Error", "server response is nil")
		return nil, diags
	}

	model := &ServerBaseModel{}

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

	if sr.Hypervisor != nil {
		model.Hypervisor = types.StringValue(*sr.Hypervisor)
	} else {
		model.Hypervisor = types.StringNull()
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

	// Objetos embebidos obligatorios
	datacenterObj, datacenterDiags := NewBaseDatacenterObject(sr.Datacenter)
	diags.Append(datacenterDiags...)
	if !datacenterDiags.HasError() {
		model.Datacenter = datacenterObj
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

	// Objetos embebidos opcionales
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

func newServerDetailModelFromResponse(ctx context.Context, sr *ServerDetailResponse) (*ServerDetailModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	baseModel, baseDiags := newServerBaseModelFromResponse(ctx, &sr.ServerBaseResponse)
	diags.Append(baseDiags...)
	if baseDiags.HasError() {
		return nil, diags
	}

	model := &ServerDetailModel{
		ServerBaseModel: *baseModel,
	}

	statusObj, statusDiags := server.NewStatusDetailObject(sr.Status)
	diags.Append(statusDiags...)
	if !statusDiags.HasError() {
		model.Status = statusObj
	}

	if sr.RecoveryMode != nil {
		model.RecoveryMode = types.BoolValue(*sr.RecoveryMode)
	} else {
		model.RecoveryMode = types.BoolNull()
	}

	if sr.RecoveryImageOS != nil {
		model.RecoveryImageOS = types.StringValue(*sr.RecoveryImageOS)
	} else {
		model.RecoveryImageOS = types.StringNull()
	}

	if sr.RecoveryUser != nil {
		model.RecoveryUser = types.StringValue(*sr.RecoveryUser)
	} else {
		model.RecoveryUser = types.StringNull()
	}

	if sr.RecoveryPassword != nil {
		model.RecoveryPassword = types.StringValue(*sr.RecoveryPassword)
	} else {
		model.RecoveryPassword = types.StringNull()
	}

	return model, diags
}

func NewServerDetailModel(ctx context.Context, sr *ServerDetailResponse) (*ServerDetailModel, diag.Diagnostics) {
	return newServerDetailModelFromResponse(ctx, sr)
}

func NewServerResourceModelFromCreate(ctx context.Context, sr *ServerDetailResponse, plan *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	baseModel, baseDiags := newServerDetailModelFromResponse(ctx, sr)
	if baseDiags.HasError() {
		diags.Append(baseDiags...)
		return nil, diags
	}

	model := &ServerResourceModel{
		ServerDetailModel: *baseModel,
	}

	// Siempre preservar los campos de configuración del plan
	model.ApplianceID = plan.ApplianceID
	model.DatacenterID = plan.DatacenterID

	// Hardware - preservar del plan si todos los campos son conocidos
	if !plan.Hardware.IsNull() && !plan.Hardware.IsUnknown() {
		hardwareAttrs := plan.Hardware.Attributes()

		allHardwareFieldsKnown := true

		if vcore, exists := hardwareAttrs["vcore"]; exists && vcore.IsUnknown() {
			allHardwareFieldsKnown = false
		}
		if ram, exists := hardwareAttrs["ram"]; exists && ram.IsUnknown() {
			allHardwareFieldsKnown = false
		}
		if hdds, exists := hardwareAttrs["hdds"]; exists && hdds.IsUnknown() {
			allHardwareFieldsKnown = false
		}
		if cores, exists := hardwareAttrs["cores_per_processor"]; exists && cores.IsUnknown() {
			allHardwareFieldsKnown = false
		}
		if fixedInstanceSizeID, exists := hardwareAttrs["fixed_instance_size_id"]; exists && fixedInstanceSizeID.IsUnknown() {
			allHardwareFieldsKnown = false
		}

		if allHardwareFieldsKnown {
			model.Hardware = plan.Hardware
		}
	}

	// Campos de configuración - preservar del plan o usar defaults
	if !plan.Password.IsUnknown() {
		model.Password = plan.Password
	} else {
		model.Password = types.StringNull()
	}

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

	// Campos opcionales de configuración
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

	if !plan.SiteID.IsUnknown() {
		model.SiteID = plan.SiteID
	} else {
		model.SiteID = types.StringNull()
	}

	// Arrays
	if !plan.PublicKey.IsUnknown() {
		model.PublicKey = plan.PublicKey
	} else {
		model.PublicKey = types.ListNull(types.StringType)
	}

	if !plan.SSHIDs.IsUnknown() {
		model.SSHIDs = plan.SSHIDs
	} else {
		model.SSHIDs = types.ListNull(types.StringType)
	}

	return model, diags
}

func NewServerResourceModelFromRead(ctx context.Context, sr *ServerDetailResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	tflog.Info(ctx, "🔍 NewServerResourceModelFromRead CALLED")
	model := *currentState

	model.Name = types.StringValue(sr.Name)
	if sr.Description != nil {
		model.Description = types.StringValue(*sr.Description)
	} else {
		model.Description = types.StringNull()
	}

	statusObj, statusDiags := server.NewStatusDetailObject(sr.Status)
	diags.Append(statusDiags...)
	if !statusDiags.HasError() {
		model.Status = statusObj
	}

	if sr.FirstPassword != nil {
		model.FirstPassword = types.StringValue(*sr.FirstPassword)
	}

	model.Hostname = types.StringValue(sr.Hostname)

	if sr.CloudPanelID != nil {
		model.CloudPanelID = types.StringValue(*sr.CloudPanelID)
	}

	if sr.Hypervisor != nil {
		model.Hypervisor = types.StringValue(*sr.Hypervisor)
	}

	if sr.Alerts != nil {
		alertsObj, alertsDiags := server.NewAlertsObject(sr.Alerts)
		diags.Append(alertsDiags...)
		if !alertsDiags.HasError() {
			model.Alerts = alertsObj
		}
	}

	if sr.DVD != nil {
		dvdObj, dvdDiags := NewIdentifierObject(*sr.DVD)
		diags.Append(dvdDiags...)
		if !dvdDiags.HasError() {
			model.DVD = dvdObj
		}
	}

	if sr.MonitoringPolicy != nil {
		monitoringObj, monitoringDiags := NewIdentifierObject(*sr.MonitoringPolicy)
		diags.Append(monitoringDiags...)
		if !monitoringDiags.HasError() {
			model.MonitoringPolicy = monitoringObj
		}
	}

	if sr.Snapshot != nil {
		snapshotObj, snapshotDiags := server.NewSnapshotObject(sr.Snapshot)
		diags.Append(snapshotDiags...)
		if !snapshotDiags.HasError() {
			model.Snapshot = snapshotObj
		}
	}

	if !currentState.Hardware.IsNull() {
		shouldUpdateHardware := false
		hardwareAttrs := currentState.Hardware.Attributes()

		tflog.Info(ctx, "🔍 Starting hardware drift detection")

		if fixedID, ok := hardwareAttrs["fixed_instance_size_id"]; ok {
			currentFixed := fixedID.(types.String).ValueString()
			apiFixed := ""
			if sr.Hardware.FixedInstanceSizeID != nil {
				apiFixed = *sr.Hardware.FixedInstanceSizeID
			}
			if currentFixed != apiFixed {
				tflog.Warn(ctx, fmt.Sprintf("HARDWARE CHANGE DETECTED in fixed_instance_size_id: state='%s' vs api='%s'", currentFixed, apiFixed))
				shouldUpdateHardware = true
			}
		}

		if modelID, ok := hardwareAttrs["baremetal_model_id"]; ok {
			currentModel := modelID.(types.String).ValueString()
			apiModel := ""
			if sr.Hardware.BaremetalModelID != nil {
				apiModel = *sr.Hardware.BaremetalModelID
			}
			if currentModel != apiModel {
				tflog.Warn(ctx, fmt.Sprintf("HARDWARE CHANGE DETECTED in baremetal_model_id: state='%s' vs api='%s'", currentModel, apiModel))
				shouldUpdateHardware = true
			}
		}

		if ram, ok := hardwareAttrs["ram"]; ok {
			currentRAM := ram.(types.Int64).ValueInt64()
			apiRAM := int64(sr.Hardware.RAM)
			if currentRAM != apiRAM {
				tflog.Warn(ctx, fmt.Sprintf("HARDWARE CHANGE DETECTED in ram: state='%d' vs api='%d'", currentRAM, apiRAM))
				shouldUpdateHardware = true
			}
		}

		if vcore, ok := hardwareAttrs["vcore"]; ok {
			currentVCore := vcore.(types.Int64).ValueInt64()
			apiVCore := int64(sr.Hardware.VCore)
			tflog.Info(ctx, fmt.Sprintf("🔍 HARDWARE COMPARISON: state vcore=%d, api vcore=%d", currentVCore, apiVCore))
			if currentVCore != apiVCore {
				// ⚠️ CONOCEMOS que hay inconsistencia en API detail endpoint
				tflog.Debug(ctx, fmt.Sprintf("API INCONSISTENCY in vcore: state='%d' vs api='%d' (preserving state value)", currentVCore, apiVCore))
				// NO marcamos shouldUpdateHardware = true para este campo
			}
		}

		if cores, ok := hardwareAttrs["cores_per_processor"]; ok {
			currentCores := cores.(types.Int64).ValueInt64()
			apiCores := int64(sr.Hardware.CoresPerProcessor)
			if currentCores != apiCores {
				tflog.Warn(ctx, fmt.Sprintf("HARDWARE CHANGE DETECTED in cores_per_processor: state='%d' vs api='%d'", currentCores, apiCores))
				shouldUpdateHardware = true
			}
		}

		if hdds, ok := hardwareAttrs["hdds"]; ok {
			currentHDDs := hdds.(types.List)
			if len(currentHDDs.Elements()) != len(sr.Hardware.HDDs) {
				tflog.Warn(ctx, fmt.Sprintf("HARDWARE CHANGE DETECTED in hdds count: state='%d' vs api='%d'", len(currentHDDs.Elements()), len(sr.Hardware.HDDs)))
				shouldUpdateHardware = true
			}
		}

		if shouldUpdateHardware {
			tflog.Info(ctx, "🔄 Updating hardware due to legitimate changes detected")
			hardwareObj, hardwareDiags := server.NewHardwareObject(sr.Hardware)
			diags.Append(hardwareDiags...)
			if !hardwareDiags.HasError() {
				model.Hardware = hardwareObj
			}
		} else {
			tflog.Info(ctx, "✅ Preserving hardware state (no legitimate changes detected)")
		}
	} else {
		// Primera vez - usar hardware de API
		tflog.Info(ctx, "🆕 First time - using hardware from API")
		hardwareObj, hardwareDiags := server.NewHardwareObject(sr.Hardware)
		diags.Append(hardwareDiags...)
		if !hardwareDiags.HasError() {
			model.Hardware = hardwareObj
		}
	}

	tflog.Info(ctx, "🔍 NewServerResourceModelFromRead COMPLETED - hardware preserved from state")
	return &model, diags
}

func NewServerResourceModelFromUpdate(_ context.Context, sr *ServerBaseResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := *currentState

	model.Name = types.StringValue(sr.Name)

	if sr.Description != nil && *sr.Description != "" {
		model.Description = types.StringValue(*sr.Description)
	} else {
		model.Description = types.StringNull()
	}

	return &model, diags
}

func NewServerResourceModelFromAPI(ctx context.Context, sr *ServerDetailResponse) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	detailModel, baseDiags := newServerDetailModelFromResponse(ctx, sr)
	if baseDiags.HasError() {
		diags.Append(baseDiags...)
		return nil, diags
	}

	model := &ServerResourceModel{
		ServerDetailModel: *detailModel,
	}

	// Campos de configuración requeridos - se asignarán en Create/Update
	model.ApplianceID = types.StringNull()
	model.DatacenterID = types.StringNull()

	// Set default values para TODOS los campos que están en ServerResourceModel
	model.Password = types.StringNull()
	model.PowerOn = types.BoolValue(true)
	model.FirewallPolicyID = types.StringNull()
	model.IPID = types.StringNull()
	model.LoadBalancerID = types.StringNull()
	model.MonitoringPolicyID = types.StringNull()
	model.PrivateNetworkID = types.StringNull()
	model.UserData = types.StringNull()
	model.UserDataContentType = types.StringNull()
	model.PublicConnectionSpeed = types.Float64Null()
	model.PrivateConnectionSpeed = types.Float64Null()
	model.BondingAllowed = types.BoolValue(false)
	model.InstallBackupAgent = types.BoolValue(false)
	model.BackupAgentTenantData = types.StringNull()
	model.BackupSize = types.StringNull()
	model.AvailabilityZoneID = types.StringNull()
	model.SiteID = types.StringNull()
	model.DisabledPassword = types.BoolValue(false)

	// Arrays
	model.PublicKey = types.ListNull(types.StringType)
	model.SSHIDs = types.ListNull(types.StringType)

	return model, diags
}

func (s *ServerResourceModel) ToCreateRequest() ServerCreateRequest {
	req := ServerCreateRequest{
		Name:         s.Name.ValueString(),
		ServerType:   "baremetal",
		ApplianceID:  s.ApplianceID.ValueString(),
		DatacenterID: s.DatacenterID.ValueString(),
		Hardware:     server.HardwareCreateRequestFromModel(s.Hardware),
	}

	req.SSHPassword = s.SSHPassword.ValueBool()
	req.RSAKey = s.RSAKey.ValueBool()
	req.InstallBackupAgent = s.InstallBackupAgent.ValueBool()

	if !s.PowerOn.IsNull() {
		req.PowerOn = s.PowerOn.ValueBool()
	} else {
		req.PowerOn = true
	}

	// Optional fields with pointers
	helper.AssignStringPtr(&req.Description, s.Description)
	helper.AssignStringPtr(&req.Password, s.Password)
	helper.AssignStringPtr(&req.FirewallPolicyID, s.FirewallPolicyID)
	helper.AssignStringPtr(&req.IPID, s.IPID)
	helper.AssignStringPtr(&req.LoadBalancerID, s.LoadBalancerID)
	helper.AssignStringPtr(&req.MonitoringPolicyID, s.MonitoringPolicyID)
	helper.AssignStringPtr(&req.PrivateNetworkID, s.PrivateNetworkID)
	helper.AssignStringPtr(&req.UserData, s.UserData)
	helper.AssignStringPtr(&req.UserDataContentType, s.UserDataContentType)
	helper.AssignStringPtr(&req.SiteID, s.SiteID)

	helper.AssignFloatPtr(&req.PublicConnectionSpeed, s.PublicConnectionSpeed)
	helper.AssignFloatPtr(&req.PrivateConnectionSpeed, s.PrivateConnectionSpeed)

	helper.AssignBoolPtr(&req.BondingAllowed, s.BondingAllowed)
	helper.AssignBoolPtr(&req.DisabledPassword, s.DisabledPassword)

	helper.AssignStringDirect(req.BackupAgentTenantData, s.BackupAgentTenantData)
	helper.AssignStringDirect(req.BackupSize, s.BackupSize)
	helper.AssignStringDirect(req.AvailabilityZoneID, s.AvailabilityZoneID)

	if !s.PublicKey.IsNull() && !s.PublicKey.IsUnknown() && len(s.PublicKey.Elements()) > 0 {
		publicKeys := make([]string, 0, len(s.PublicKey.Elements()))
		for _, key := range s.PublicKey.Elements() {
			if strVal, ok := key.(types.String); ok {
				publicKeys = append(publicKeys, strVal.ValueString())
			}
		}
		req.PublicKey = publicKeys
	}

	if !s.SSHIDs.IsNull() && !s.SSHIDs.IsUnknown() && len(s.SSHIDs.Elements()) > 0 {
		sshIds := make([]string, 0, len(s.SSHIDs.Elements()))
		for _, id := range s.SSHIDs.Elements() {
			if strVal, ok := id.(types.String); ok {
				sshIds = append(sshIds, strVal.ValueString())
			}
		}
		req.SSHIDs = sshIds
	}

	return req
}

func (s *ServerResourceModel) ToUpdateRequest() ServerUpdateRequest {
	req := ServerUpdateRequest{}

	helper.AssignStringPtr(&req.Name, s.Name)
	helper.AssignStringPtr(&req.Description, s.Description)

	return req
}

func serverBaseModelObjectType() types.ObjectType {
	return types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"id":                types.StringType,
			"name":              types.StringType,
			"description":       types.StringType,
			"datacenter":        baseDatacenterObjectType(),
			"creation_date":     types.StringType,
			"first_password":    types.StringType,
			"managed":           types.BoolType,
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
				Attributes:  server.StatusDetailDataSourceSchema(),
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
			"recovery_mode": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether server is in recovery mode",
			},
			"recovery_image_os": schema.StringAttribute{
				Computed:    true,
				Description: "Recovery image OS",
			},
			"recovery_user": schema.StringAttribute{
				Computed:    true,
				Description: "Recovery user",
			},
			"recovery_password": schema.StringAttribute{
				Computed:    true,
				Description: "Recovery password",
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
				Optional:    true,
				Computed:    true,
				Description: "Server name",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(util.MaxNameLength),
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server description",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"datacenter": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseDatacenterResourceAttributes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"creation_date": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_password": rschema.StringAttribute{
				Computed: true,
				//Sensitive: true,
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
			"ips": rschema.ListNestedAttribute{
				Computed: true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersIPResourceSchema(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"ssh_password": rschema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Description: "Whether SSH password authentication is enabled",
			},
			"image": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierResourceAttributes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"hardware": rschema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server hardware configuration",
				Attributes:  server.HardwareResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"dvd": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierResourceAttributes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"alerts": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server alerts",
				Attributes:  server.AlertsResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"monitoring_policy": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: BaseIdentifierResourceAttributes(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
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
				Computed:    true,
				Description: "Server hypervisor type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"hostname": rschema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_speed": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Connection speed configuration",
				Attributes:  server.ConnectionSpeedResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"redundancy": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.RedundancyResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"rsa_key": rschema.BoolAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"snapshot": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.SnapshotResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"private_networks": rschema.ListNestedAttribute{
				Computed: true,
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersPrivateNetworksResourceSchema(),
				},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			// ServerDetailModel
			"status": rschema.SingleNestedAttribute{
				Computed:   true,
				Attributes: server.StatusResourceSchema(),
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"recovery_mode": rschema.BoolAttribute{
				Computed:    true,
				Description: "Whether server is in recovery mode",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"recovery_image_os": rschema.StringAttribute{
				Computed:    true,
				Description: "Recovery image OS",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"recovery_user": rschema.StringAttribute{
				Computed:    true,
				Description: "Recovery user",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"recovery_password": rschema.StringAttribute{
				Computed:    true,
				Description: "Recovery password",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			// Fields from ServerResourceModel
			"appliance_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Appliance identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid appliance_id",
					),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Datacenter identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter_id",
					),
				},
			},
			"password": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server password",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
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
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"firewall_policy_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Firewall policy identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid firewall_policy_id",
					),
				},
			},
			"ip_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "IP identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid ip_id",
					),
				},
			},
			"load_balancer_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Load balancer identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid load_balancer_id",
					),
				},
			},
			"monitoring_policy_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Monitoring policy identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid monitoring_policy_id",
					),
				},
			},
			"private_network_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Private network identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid private_network_id",
					),
				},
			},
			"user_data": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "User data script",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"user_data_content_type": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "User data content type",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_connection_speed": rschema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Public connection speed",
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"private_connection_speed": rschema.Float64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Private connection speed",
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"bonding_allowed": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether bonding is allowed",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"install_backup_agent": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to install backup agent",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_agent_tenant_data": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup agent tenant data",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"backup_size": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Backup size",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"availability_zone_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Availability zone identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"public_key": rschema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "List of SSH Key IDs to be copied in the server",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"site_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Site identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid site_id",
					),
				},
			},
			"ssh_ids": rschema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Description: "SSH key identifiers",
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"disabled_password": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Default:     booldefault.StaticBool(false),
				Description: "Whether to disable password authentication",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}
