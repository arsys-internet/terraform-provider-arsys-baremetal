package models

import (
	"context"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/models/server"
	"terraform-provider-arsys-baremetal/internal/util"
	"terraform-provider-arsys-baremetal/internal/util/helper"

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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type ServerBaseModel struct {
	Id               types.String `tfsdk:"id"`
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
	CloudPanelId     types.String `tfsdk:"cloudpanel_id"`
	ServerType       types.String `tfsdk:"server_type"`
	Hostname         types.String `tfsdk:"hostname"`
	ConnectionSpeed  types.Object `tfsdk:"connection_speed"`
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

	ApplianceId        types.String `tfsdk:"appliance_id"`
	DatacenterId       types.String `tfsdk:"datacenter_id"`
	Password           types.String `tfsdk:"password"`
	PowerOn            types.Bool   `tfsdk:"power_on"`
	FirewallPolicyId   types.String `tfsdk:"firewall_policy_id"`
	LoadBalancerId     types.String `tfsdk:"load_balancer_id"`
	MonitoringPolicyId types.String `tfsdk:"monitoring_policy_id"`
	InstallBackupAgent types.Bool   `tfsdk:"install_backup_agent"`
	AvailabilityZoneId types.String `tfsdk:"availability_zone_id"`
	PublicKey          types.List   `tfsdk:"public_key"`
}

type ServerBaseResponse struct {
	Id               string                                 `json:"id"`
	Name             string                                 `json:"name"`
	Description      *string                                `json:"description"`
	Datacenter       BaseDatacenterResponse                 `json:"datacenter"`
	CreationDate     string                                 `json:"creation_date"`
	FirstPassword    *string                                `json:"first_password"`
	Managed          bool                                   `json:"managed"`
	IPs              []server.ServersIPResponse             `json:"ips"`
	SSHPassword      bool                                   `json:"ssh_password"`
	Image            IdentifierResponse                     `json:"image"`
	Hardware         server.HardwareResponse                `json:"hardware"`
	DVD              *IdentifierResponse                    `json:"dvd"`
	Alerts           *server.AlertResponse                  `json:"alerts"`
	MonitoringPolicy *IdentifierResponse                    `json:"monitoring_policy"`
	CloudPanelId     *string                                `json:"cloudpanel_id"`
	ServerType       string                                 `json:"server_type"`
	Hostname         string                                 `json:"hostname"`
	ConnectionSpeed  *server.ConnectionSpeedResponse        `json:"connection_speed"`
	Redundancy       *server.RedundancyResponse             `json:"redundancy"`
	RSAKey           interface{}                            `json:"rsa_key"`
	Snapshot         *server.SnapshotResponse               `json:"snapshot"`
	PrivateNetworks  []server.ServersPrivateNetworkResponse `json:"private_networks"`
}

type ServerDetailResponse struct {
	ServerBaseResponse
	Status           server.StatusDetailResponse `json:"status"`
	RecoveryMode     bool                        `json:"recovery_mode"`
	RecoveryImageOS  *string                     `json:"recovery_image_os"`
	RecoveryUser     *string                     `json:"recovery_user"`
	RecoveryPassword *string                     `json:"recovery_password"`
}

type ServerCreateRequest struct {
	Name         string                       `json:"name"`
	ServerType   string                       `json:"server_type"`
	ApplianceId  string                       `json:"appliance_id"`
	DatacenterId string                       `json:"datacenter_id"`
	Hardware     server.HardwareCreateRequest `json:"hardware"`

	SSHPassword        bool `json:"ssh_password"`
	PowerOn            bool `json:"power_on"`
	RSAKey             bool `json:"rsa_key"`
	InstallBackupAgent bool `json:"install_backup_agent"`

	Description        *string `json:"description,omitempty"`
	Password           *string `json:"password,omitempty"`
	FirewallPolicyId   *string `json:"firewall_policy_id,omitempty"`
	LoadBalancerId     *string `json:"load_balancer_id,omitempty"`
	MonitoringPolicyId *string `json:"monitoring_policy_id,omitempty"`
	AvailabilityZoneId *string `json:"availability_zone_id,omitempty"`
	PublicKey          []string `json:"public_key,omitempty"`
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

	model.Id = types.StringValue(sr.Id)
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

	if sr.CloudPanelId != nil {
		model.CloudPanelId = types.StringValue(*sr.CloudPanelId)
	} else {
		model.CloudPanelId = types.StringNull()
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

	model.RecoveryMode = types.BoolValue(sr.RecoveryMode)

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

	model.ApplianceId = plan.ApplianceId
	model.DatacenterId = plan.DatacenterId

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
		if fixedInstanceSizeId, exists := hardwareAttrs["fixed_instance_size_id"]; exists && fixedInstanceSizeId.IsUnknown() {
			allHardwareFieldsKnown = false
		}

		if allHardwareFieldsKnown {
			model.Hardware = plan.Hardware
		}
	}

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

	if !plan.InstallBackupAgent.IsUnknown() {
		model.InstallBackupAgent = plan.InstallBackupAgent
	} else {
		model.InstallBackupAgent = types.BoolValue(false)
	}

	if !plan.FirewallPolicyId.IsUnknown() {
		model.FirewallPolicyId = plan.FirewallPolicyId
	} else {
		model.FirewallPolicyId = types.StringNull()
	}

	if !plan.LoadBalancerId.IsUnknown() {
		model.LoadBalancerId = plan.LoadBalancerId
	} else {
		model.LoadBalancerId = types.StringNull()
	}

	if !plan.MonitoringPolicyId.IsUnknown() {
		model.MonitoringPolicyId = plan.MonitoringPolicyId
	} else {
		model.MonitoringPolicyId = types.StringNull()
	}

	if !plan.AvailabilityZoneId.IsUnknown() {
		model.AvailabilityZoneId = plan.AvailabilityZoneId
	} else {
		model.AvailabilityZoneId = types.StringNull()
	}

	if !plan.PublicKey.IsUnknown() {
		model.PublicKey = plan.PublicKey
	} else {
		model.PublicKey = types.ListNull(types.StringType)
	}

	return model, diags
}

func NewServerResourceModelFromRead(_ context.Context, sr *ServerDetailResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

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

	if sr.CloudPanelId != nil {
		model.CloudPanelId = types.StringValue(*sr.CloudPanelId)
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

	shouldUpdateHardware := currentState.Hardware.IsNull() ||
		server.NeedsHardwareUpdate(currentState.Hardware.Attributes(), sr.Hardware)

	if shouldUpdateHardware {
		hardwareObj, hardwareDiags := server.NewHardwareObject(sr.Hardware)
		diags.Append(hardwareDiags...)
		if !hardwareDiags.HasError() {
			model.Hardware = hardwareObj
		}
	}

	return &model, diags
}

func NewServerResourceModelFromUpdate(_ context.Context, sr *ServerBaseResponse, currentState *ServerResourceModel) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	model := *currentState

	model.Name = types.StringValue(sr.Name)

	helper.StringPtrToTypesStringWithNullEmpty(&model.Description, sr.Description)

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

	model.ApplianceId = types.StringNull()
	model.DatacenterId = types.StringNull()

	model.Password = types.StringNull()
	model.PowerOn = types.BoolValue(true)
	model.FirewallPolicyId = types.StringNull()
	model.LoadBalancerId = types.StringNull()
	model.MonitoringPolicyId = types.StringNull()
	model.InstallBackupAgent = types.BoolValue(false)
	model.AvailabilityZoneId = types.StringNull()

	return model, diags
}

func (s *ServerResourceModel) ToCreateRequest() ServerCreateRequest {
	req := ServerCreateRequest{
		Name:         s.Name.ValueString(),
		ServerType:   "baremetal",
		ApplianceId:  s.ApplianceId.ValueString(),
		DatacenterId: s.DatacenterId.ValueString(),
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

	helper.AssignStringPtr(&req.Description, s.Description)
	helper.AssignStringPtr(&req.Password, s.Password)
	helper.AssignStringPtr(&req.FirewallPolicyId, s.FirewallPolicyId)
	helper.AssignStringPtr(&req.LoadBalancerId, s.LoadBalancerId)
	helper.AssignStringPtr(&req.MonitoringPolicyId, s.MonitoringPolicyId)
	helper.AssignStringPtr(&req.AvailabilityZoneId, s.AvailabilityZoneId)

	if !s.PublicKey.IsNull() && !s.PublicKey.IsUnknown() {
		var keys []string
		s.PublicKey.ElementsAs(context.Background(), &keys, false)
		req.PublicKey = keys
	}

	return req
}

func NewServerResourceModelFromImport(ctx context.Context, sr *ServerDetailResponse) (*ServerResourceModel, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	detailModel, baseDiags := newServerDetailModelFromResponse(ctx, sr)
	if baseDiags.HasError() {
		diags.Append(baseDiags...)
		return nil, diags
	}

	model := &ServerResourceModel{
		ServerDetailModel: *detailModel,
	}

	model.ApplianceId = types.StringNull()
	model.DatacenterId = types.StringNull()
	model.Password = types.StringNull()
	model.PowerOn = types.BoolNull()
	model.FirewallPolicyId = types.StringNull()
	model.LoadBalancerId = types.StringNull()
	model.MonitoringPolicyId = types.StringNull()
	model.InstallBackupAgent = types.BoolNull()
	model.AvailabilityZoneId = types.StringNull()
	model.PublicKey = types.ListNull(types.StringType)

	return model, diags
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
						"must be a valid server ID",
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
				Computed:    true,
				Description: "Server image",
				Attributes:  BaseIdentifierAttributes(),
			},
			"hardware": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Hardware configuration",
				Attributes:  server.HardwareDataSourceSchema(),
			},
			"dvd": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "DVD image",
				Attributes:  BaseIdentifierAttributes(),
			},
			"alerts": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server alerts",
				Attributes:  server.AlertsDataSourceSchema(),
			},
			"monitoring_policy": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Monitoring policy",
				Attributes:  BaseIdentifierAttributes(),
			},
			"cloudpanel_id": schema.StringAttribute{
				Description: "CloudPanel Id",
				Computed:    true,
			},
			"server_type": schema.StringAttribute{
				Description: "Server type",
				Computed:    true,
			},
			"hostname": schema.StringAttribute{
				Description: "Hostname",
				Computed:    true,
			},
			"connection_speed": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Connection speed configuration",
				Attributes:  server.ConnectionSpeedDataSourceSchema(),
			},
			"redundancy": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Redundancy configuration",
				Attributes:  server.RedundancyDataSourceSchema(),
			},
			"rsa_key": schema.BoolAttribute{
				Computed:    true,
				Description: "Whether RSA key authentication is enabled",
			},
			"snapshot": schema.SingleNestedAttribute{
				Computed:    true,
				Description: "Snapshot configuration",
				Attributes:  server.SnapshotDataSourceSchema(),
			},
			"private_networks": schema.ListNestedAttribute{
				Computed:    true,
				Description: "Private networks configuration",
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
				Description: "Server description",
				Validators: []validator.String{
					stringvalidator.LengthAtMost(256),
				},
			},
			"datacenter": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server datacenter",
				Attributes:  BaseDatacenterResourceAttributes(),
			},
			"creation_date": rschema.StringAttribute{
				Computed:    true,
				Description: "Creation timestamp",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"first_password": rschema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "First password generated",
			},
			"managed": rschema.BoolAttribute{
				Computed:    true,
				Description: "Whether server is managed",
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ips": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "Server IP addresses",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersIPResourceSchema(),
				},
			},
			"ssh_password": rschema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Whether SSH password authentication is enabled",
			},
			"image": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server image",
				Attributes:  BaseIdentifierResourceAttributes(),
			},
			"hardware": rschema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server hardware configuration",
				Attributes:  server.HardwareResourceSchema(),
			},
			"dvd": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "DVD image",
				Attributes:  BaseIdentifierResourceAttributes(),
			},
			"alerts": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server alerts",
				Attributes:  server.AlertsResourceSchema(),
			},
			"monitoring_policy": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Monitoring policy",
				Attributes:  BaseIdentifierResourceAttributes(),
			},
			"cloudpanel_id": rschema.StringAttribute{
				Computed:    true,
				Description: "CloudPanel Id",
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
			"hostname": rschema.StringAttribute{
				Computed:    true,
				Description: "Hostname",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"connection_speed": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Connection speed configuration",
				Attributes:  server.ConnectionSpeedResourceSchema(),
			},
			"redundancy": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Redundancy configuration",
				Attributes:  server.RedundancyResourceSchema(),
			},
			"rsa_key": rschema.BoolAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Whether RSA key authentication is enabled",
			},
			"snapshot": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Snapshot configuration",
				Attributes:  server.SnapshotResourceSchema(),
			},
			"private_networks": rschema.ListNestedAttribute{
				Computed:    true,
				Description: "Private networks configuration",
				NestedObject: rschema.NestedAttributeObject{
					Attributes: server.ServersPrivateNetworksResourceSchema(),
				},
			},
			"status": rschema.SingleNestedAttribute{
				Computed:    true,
				Description: "Server status",
				Attributes:  server.StatusDetailResourceSchema(),
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
			"appliance_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Appliance identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid appliance ID",
					),
				},
			},
			"datacenter_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Datacenter identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid datacenter ID",
					),
				},
			},
			"password": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Server password",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(8),
					stringvalidator.LengthAtMost(64),
				},
			},
			"power_on": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to power on the server after creation",
			},
			"firewall_policy_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Firewall policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid firewall policy ID",
					),
				},
			},
			"load_balancer_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Load balancer identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid load balancer ID",
					),
				},
			},
			"monitoring_policy_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Monitoring policy identifier",
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(util.HexID32Pattern),
						"must be a valid monitoring policy ID",
					),
				},
			},
			"install_backup_agent": rschema.BoolAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Whether to install backup agent",
			},
			"availability_zone_id": rschema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "Availability zone identifier",
			},
			"public_key": rschema.ListAttribute{
				Optional:    true,
				Description: "List of SSH Key IDs to be copied in the server. Then you will be able to access to the server using your SSH keys.",
				ElementType: types.StringType,
			},
		},
	}
}
