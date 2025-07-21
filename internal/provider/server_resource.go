package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/server"
	"terraform-provider-arsys-baremetal/internal/util"
	"time"
)

var (
	_ resource.Resource                = &ServerResource{}
	_ resource.ResourceWithConfigure   = &ServerResource{}
	_ resource.ResourceWithImportState = &ServerResource{}
)

func NewServerResource() resource.Resource {
	return &ServerResource{}
}

type ServerResource struct {
	client *service.ApiServerService
}

func (r *ServerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (r *ServerResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.ServerResourceSchema(ctx)
}

func (r *ServerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetServerService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	serverService, ok := client.(*service.ApiServerService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = serverService
}

func (r *ServerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Missing required field",
			"'name' field is required when creating a server",
		)
	}

	if data.ApplianceID.IsNull() || data.ApplianceID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("appliance_id"),
			"Missing required field",
			"'appliance_id' field is required when creating a server",
		)
	}

	if data.DatacenterID.IsNull() || data.DatacenterID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("datacenter_id"),
			"Missing required field",
			"Either 'datacenter_id' field is required when creating a server",
		)
	}

	if data.Hardware.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("hardware"),
			"Missing required field",
			"'hardware' field is required when creating a server",
		)
	}

	hardwareAttrs := data.Hardware.Attributes()
	if baremetalModelID, ok := hardwareAttrs["baremetal_model_id"].(types.String); !ok || baremetalModelID.IsNull() || baremetalModelID.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("hardware").AtName("baremetal_model_id"),
			"Missing required field",
			"'baremetal_model_id' is required when creating a baremetal server",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest := data.ToCreateRequest()

	tflog.Info(ctx, fmt.Sprintf("Creating server: %s", createRequest.Name))

	apiResponse, err := r.client.CreateServer(&createRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			fmt.Sprintf("Could not create server: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getServerTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
		[]string{util.StateDeploying},
		[]string{util.StatePoweredOn, util.StatePoweredOff, util.StateActive},
	)

	_, diags := util.WaitForResourceState(
		ctx,
		apiResponse.ID,
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Wait for server state failed")
		return
	}
	tflog.Info(ctx, fmt.Sprintf("apiResponse: %+v", apiResponse))

	finalServer, err := r.client.GetServer(apiResponse.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting final server state",
			fmt.Sprintf("Could not get final server state: %s", err),
		)
		return
	}
	debugServerResponse(ctx, finalServer, "Create Request Final Server")

	tflog.Info(ctx, fmt.Sprintf("Create - finalServer: %+v", finalServer))

	finalModel, diags := models.NewServerResourceModelFromCreate(ctx, finalServer, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to create final resource model")
		return
	}
	debugTerraformModel(ctx, finalModel, "Create Request Final Model")

	tflog.Info(ctx, fmt.Sprintf("Created server with ID: %s", finalModel.ID.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *ServerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "🔥 READ RESOURCE FUNCTION CALLED - STARTING")
	var data models.ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading server with ID: %s", id))

	apiResponse, err := r.client.GetServer(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Server with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading server",
			fmt.Sprintf("Could not read server: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Read - apiResponse: %+v", apiResponse))

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Server with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}
	debugServerResponse(ctx, apiResponse, "Read Request Api response")

	readModel, diags := models.NewServerResourceModelFromRead(ctx, apiResponse, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	debugTerraformModel(ctx, readModel, "Read Request Read Model")
	tflog.Info(ctx, "🔥 READ RESOURCE FUNCTION COMPLETED - model preserved")

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *ServerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.ServerResourceModel
	var state models.ServerResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating server with ID: %s", id))

	hasChanges := false
	if !plan.Name.Equal(state.Name) {
		hasChanges = true
		tflog.Info(ctx, "Name changed")
	}

	if !plan.Description.Equal(state.Description) {
		hasChanges = true
		tflog.Info(ctx, "Description changed")
	}

	if !hasChanges {
		tflog.Info(ctx, "No changes detected, skipping API call")
		resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
		return
	}

	updateRequest := plan.ToUpdateRequest()
	tflog.Info(ctx, fmt.Sprintf("Update request: name=%s, description=%s",
		updateRequest.Name, updateRequest.Description))

	updatedServer, err := r.client.UpdateServer(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating server",
			fmt.Sprintf("Could not update server: %s", err),
		)
		return
	}

	finalModel, diags := models.NewServerResourceModelFromUpdate(ctx, updatedServer, &state)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, fmt.Sprintf("Successfully updated server with ID: %s", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *ServerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.ServerResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting server with ID: %s", id))

	err := r.client.DeleteServer(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Server %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting server",
			fmt.Sprintf("Could not delete server: %s", err),
		)
		return
	}

	defaultTimeout, defaultRetryInterval, defaultMinTimeout := getServerTimeout()

	waitOptions := util.NewWaitOptions(
		defaultTimeout,
		defaultRetryInterval,
		defaultMinTimeout,
		[]string{util.StateRemoving},
		[]string{util.StateDeleted},
	)

	waitOptions.IgnoreNotFoundErrors = true

	_, diags := util.WaitForResourceState(
		ctx,
		id,
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Deleted server with ID: %s", id))
}

func (r *ServerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getServerTimeout() (time.Duration, time.Duration, time.Duration) {
	timeout, err := util.GetEnvTimeValues("SERVER_DEFAULT_TIMEOUT", time.Minute)
	if err != nil {
		timeout = 40 * time.Minute
	}

	retryInterval, err := util.GetEnvTimeValues("SERVER_DEFAULT_RETRY_INTERVAL", time.Second)
	if err != nil {
		retryInterval = 30 * time.Second
	}

	minTimeout, err := util.GetEnvTimeValues("SERVER_DEFAULT_MIN_TIMEOUT", time.Second)
	if err != nil {
		minTimeout = 10 * time.Second
	}

	return timeout, retryInterval, minTimeout
}

func debugTerraformModel(ctx context.Context, model *models.ServerResourceModel, label string) {
	if model == nil {
		tflog.Info(ctx, fmt.Sprintf("📄 %s: MODEL IS NIL", label))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("📄 %s:", label))

	// ✅ Campos básicos obligatorios
	tflog.Info(ctx, fmt.Sprintf("  ID: %s (null: %t, unknown: %t)",
		model.ID.ValueString(), model.ID.IsNull(), model.ID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Name: %s (null: %t, unknown: %t)",
		model.Name.ValueString(), model.Name.IsNull(), model.Name.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Description: %s (null: %t, unknown: %t)",
		model.Description.ValueString(), model.Description.IsNull(), model.Description.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  ApplianceID: %s (null: %t, unknown: %t)",
		model.ApplianceID.ValueString(), model.ApplianceID.IsNull(), model.ApplianceID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  DatacenterID: %s (null: %t, unknown: %t)",
		model.DatacenterID.ValueString(), model.DatacenterID.IsNull(), model.DatacenterID.IsUnknown()))

	// ✅ Campos computed importantes
	tflog.Info(ctx, fmt.Sprintf("  CreationDate: %s (null: %t, unknown: %t)",
		model.CreationDate.ValueString(), model.CreationDate.IsNull(), model.CreationDate.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  FirstPassword: %s (null: %t, unknown: %t)",
		model.FirstPassword.ValueString(), model.FirstPassword.IsNull(), model.FirstPassword.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  ServerType: %s (null: %t, unknown: %t)",
		model.ServerType.ValueString(), model.ServerType.IsNull(), model.ServerType.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Hostname: %s (null: %t, unknown: %t)",
		model.Hostname.ValueString(), model.Hostname.IsNull(), model.Hostname.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  CloudPanelID: %s (null: %t, unknown: %t)",
		model.CloudPanelID.ValueString(), model.CloudPanelID.IsNull(), model.CloudPanelID.IsUnknown()))
	// ✅ Campos booleanos
	tflog.Info(ctx, fmt.Sprintf("  PowerOn: %t (null: %t, unknown: %t)",
		model.PowerOn.ValueBool(), model.PowerOn.IsNull(), model.PowerOn.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  SSHPassword: %t (null: %t, unknown: %t)",
		model.SSHPassword.ValueBool(), model.SSHPassword.IsNull(), model.SSHPassword.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Managed: %t (null: %t, unknown: %t)",
		model.Managed.ValueBool(), model.Managed.IsNull(), model.Managed.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  RSAKey: %t (null: %t, unknown: %t)",
		model.RSAKey.ValueBool(), model.RSAKey.IsNull(), model.RSAKey.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  BondingAllowed: %t (null: %t, unknown: %t)",
		model.BondingAllowed.ValueBool(), model.BondingAllowed.IsNull(), model.BondingAllowed.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  InstallBackupAgent: %t (null: %t, unknown: %t)",
		model.InstallBackupAgent.ValueBool(), model.InstallBackupAgent.IsNull(), model.InstallBackupAgent.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  DisabledPassword: %t (null: %t, unknown: %t)",
		model.DisabledPassword.ValueBool(), model.DisabledPassword.IsNull(), model.DisabledPassword.IsUnknown()))

	// ✅ Campos opcionales de configuración
	tflog.Info(ctx, fmt.Sprintf("  Password: [HIDDEN] (null: %t, unknown: %t)",
		model.Password.IsNull(), model.Password.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  FirewallPolicyID: %s (null: %t, unknown: %t)",
		model.FirewallPolicyID.ValueString(), model.FirewallPolicyID.IsNull(), model.FirewallPolicyID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  IPID: %s (null: %t, unknown: %t)",
		model.IPID.ValueString(), model.IPID.IsNull(), model.IPID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  LoadBalancerID: %s (null: %t, unknown: %t)",
		model.LoadBalancerID.ValueString(), model.LoadBalancerID.IsNull(), model.LoadBalancerID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  MonitoringPolicyID: %s (null: %t, unknown: %t)",
		model.MonitoringPolicyID.ValueString(), model.MonitoringPolicyID.IsNull(), model.MonitoringPolicyID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  PrivateNetworkID: %s (null: %t, unknown: %t)",
		model.PrivateNetworkID.ValueString(), model.PrivateNetworkID.IsNull(), model.PrivateNetworkID.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  PublicKey: [HIDDEN] (null: %t, unknown: %t)",
		model.PublicKey.IsNull(), model.PublicKey.IsUnknown()))

	// ✅ Campos numéricos
	tflog.Info(ctx, fmt.Sprintf("  PublicConnectionSpeed: %f (null: %t, unknown: %t)",
		model.PublicConnectionSpeed.ValueFloat64(), model.PublicConnectionSpeed.IsNull(), model.PublicConnectionSpeed.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  PrivateConnectionSpeed: %f (null: %t, unknown: %t)",
		model.PrivateConnectionSpeed.ValueFloat64(), model.PrivateConnectionSpeed.IsNull(), model.PrivateConnectionSpeed.IsUnknown()))

	// ✅ Hardware (objeto complex)
	tflog.Info(ctx, fmt.Sprintf("  Hardware: (null: %t, unknown: %t)",
		model.Hardware.IsNull(), model.Hardware.IsUnknown()))

	if !model.Hardware.IsNull() && !model.Hardware.IsUnknown() {
		hardwareAttrs := model.Hardware.Attributes()

		if baremetalID, ok := hardwareAttrs["baremetal_model_id"]; ok {
			tflog.Info(ctx, fmt.Sprintf("    BaremetalModelID: %s (null: %t, unknown: %t)",
				baremetalID.(types.String).ValueString(),
				baremetalID.(types.String).IsNull(),
				baremetalID.(types.String).IsUnknown()))
		}

		if vcore, ok := hardwareAttrs["vcore"]; ok {
			tflog.Info(ctx, fmt.Sprintf("    VCore: %d (null: %t, unknown: %t)",
				vcore.(types.Int64).ValueInt64(),
				vcore.(types.Int64).IsNull(),
				vcore.(types.Int64).IsUnknown()))
		}

		if ram, ok := hardwareAttrs["ram"]; ok {
			tflog.Info(ctx, fmt.Sprintf("    RAM: %d (null: %t, unknown: %t)",
				ram.(types.Int64).ValueInt64(),
				ram.(types.Int64).IsNull(),
				ram.(types.Int64).IsUnknown()))
		}

		if cores, ok := hardwareAttrs["cores_per_processor"]; ok {
			tflog.Info(ctx, fmt.Sprintf("    CoresPerProcessor: %d (null: %t, unknown: %t)",
				cores.(types.Int64).ValueInt64(),
				cores.(types.Int64).IsNull(),
				cores.(types.Int64).IsUnknown()))
		}

		if hdds, ok := hardwareAttrs["hdds"]; ok {
			tflog.Info(ctx, fmt.Sprintf("    HDDs: (null: %t, unknown: %t, elements: %d)",
				hdds.(types.List).IsNull(),
				hdds.(types.List).IsUnknown(),
				len(hdds.(types.List).Elements())))
		}

		if fixedInstanceID, ok := hardwareAttrs["fixed_instance_size_id"]; ok {
			tflog.Info(ctx, fmt.Sprintf("    FixedInstanceSizeID: %s (null: %t, unknown: %t)",
				fixedInstanceID.(types.String).ValueString(),
				fixedInstanceID.(types.String).IsNull(),
				fixedInstanceID.(types.String).IsUnknown()))
		}
	} else {
		tflog.Info(ctx, "    Hardware is null or unknown - no attributes available")
	}

	// ✅ Objetos nested computed
	tflog.Info(ctx, fmt.Sprintf("  Status: (null: %t, unknown: %t)",
		model.Status.IsNull(), model.Status.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  ConnectionSpeed: (null: %t, unknown: %t)",
		model.ConnectionSpeed.IsNull(), model.ConnectionSpeed.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Datacenter: (null: %t, unknown: %t)",
		model.Datacenter.IsNull(), model.Datacenter.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Image: (null: %t, unknown: %t)",
		model.Image.IsNull(), model.Image.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  DVD: (null: %t, unknown: %t)",
		model.DVD.IsNull(), model.DVD.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Alerts: (null: %t, unknown: %t)",
		model.Alerts.IsNull(), model.Alerts.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  MonitoringPolicy: (null: %t, unknown: %t)",
		model.MonitoringPolicy.IsNull(), model.MonitoringPolicy.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Redundancy: (null: %t, unknown: %t)",
		model.Redundancy.IsNull(), model.Redundancy.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  Snapshot: (null: %t, unknown: %t)",
		model.Snapshot.IsNull(), model.Snapshot.IsUnknown()))

	// ✅ Listas computed
	tflog.Info(ctx, fmt.Sprintf("  IPs: (null: %t, unknown: %t, elements: %d)",
		model.IPs.IsNull(), model.IPs.IsUnknown(), len(model.IPs.Elements())))
	tflog.Info(ctx, fmt.Sprintf("  PrivateNetworks: (null: %t, unknown: %t, elements: %d)",
		model.PrivateNetworks.IsNull(), model.PrivateNetworks.IsUnknown(), len(model.PrivateNetworks.Elements())))
	tflog.Info(ctx, fmt.Sprintf("  SSHIDs: (null: %t, unknown: %t, elements: %d)",
		model.SSHIDs.IsNull(), model.SSHIDs.IsUnknown(), len(model.SSHIDs.Elements())))

	// ✅ Campos de backup y configuración avanzada
	tflog.Info(ctx, fmt.Sprintf("  BackupAgentTenantData: %s (null: %t, unknown: %t)",
		model.BackupAgentTenantData.ValueString(), model.BackupAgentTenantData.IsNull(), model.BackupAgentTenantData.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  BackupSize: %s (null: %t, unknown: %t)",
		model.BackupSize.ValueString(), model.BackupSize.IsNull(), model.BackupSize.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  AvailabilityZoneID: %s (null: %t, unknown: %t)",
		model.AvailabilityZoneID.ValueString(), model.AvailabilityZoneID.IsNull(), model.AvailabilityZoneID.IsUnknown()))

	// ✅ User data y contenido
	tflog.Info(ctx, fmt.Sprintf("  UserData: [HIDDEN] (null: %t, unknown: %t)",
		model.UserData.IsNull(), model.UserData.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  UserDataContentType: %s (null: %t, unknown: %t)",
		model.UserDataContentType.ValueString(), model.UserDataContentType.IsNull(), model.UserDataContentType.IsUnknown()))

	tflog.Info(ctx, fmt.Sprintf("  SiteID: %s (null: %t, unknown: %t)",
		model.SiteID.ValueString(), model.SiteID.IsNull(), model.SiteID.IsUnknown()))

	// Agregar campos de recuperación después de UserDataContentType:
	tflog.Info(ctx, fmt.Sprintf("  RecoveryMode: %t (null: %t, unknown: %t)",
		model.RecoveryMode.ValueBool(), model.RecoveryMode.IsNull(), model.RecoveryMode.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  RecoveryImageOS: %s (null: %t, unknown: %t)",
		model.RecoveryImageOS.ValueString(), model.RecoveryImageOS.IsNull(), model.RecoveryImageOS.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  RecoveryUser: %s (null: %t, unknown: %t)",
		model.RecoveryUser.ValueString(), model.RecoveryUser.IsNull(), model.RecoveryUser.IsUnknown()))
	tflog.Info(ctx, fmt.Sprintf("  RecoveryPassword: [HIDDEN] (null: %t, unknown: %t)",
		model.RecoveryPassword.IsNull(), model.RecoveryPassword.IsUnknown()))
}

func debugServerResponse(ctx context.Context, sr *models.ServerDetailResponse, label string) {
	if sr == nil {
		tflog.Info(ctx, fmt.Sprintf("🔧 %s: SERVER RESPONSE IS NIL", label))
		return
	}

	tflog.Info(ctx, fmt.Sprintf("🔧 %s:", label))

	// ✅ Campos básicos obligatorios
	tflog.Info(ctx, fmt.Sprintf("  ID: %s", sr.ID))
	tflog.Info(ctx, fmt.Sprintf("  Name: %s", sr.Name))
	tflog.Info(ctx, fmt.Sprintf("  CreationDate: %s", sr.CreationDate))
	tflog.Info(ctx, fmt.Sprintf("  Managed: %t", sr.Managed))
	tflog.Info(ctx, fmt.Sprintf("  SSHPassword: %t", sr.SSHPassword))
	tflog.Info(ctx, fmt.Sprintf("  ServerType: %s", sr.ServerType))
	tflog.Info(ctx, fmt.Sprintf("  Hostname: %s", sr.Hostname))

	if sr.RecoveryMode != nil {
		tflog.Info(ctx, fmt.Sprintf("  RecoveryMode: %t", *sr.RecoveryMode))
	} else {
		tflog.Info(ctx, "  RecoveryMode: <nil>")
	}

	if sr.RecoveryImageOS != nil {
		tflog.Info(ctx, fmt.Sprintf("  RecoveryImageOS: %s", *sr.RecoveryImageOS))
	} else {
		tflog.Info(ctx, "  RecoveryImageOS: <nil>")
	}

	if sr.RecoveryUser != nil {
		tflog.Info(ctx, fmt.Sprintf("  RecoveryUser: %s", *sr.RecoveryUser))
	} else {
		tflog.Info(ctx, "  RecoveryUser: <nil>")
	}

	if sr.RecoveryPassword != nil {
		tflog.Info(ctx, "  RecoveryPassword: [HIDDEN]")
	} else {
		tflog.Info(ctx, "  RecoveryPassword: <nil>")
	}

	// ✅ Campos opcionales (punteros)
	if sr.Description != nil {
		tflog.Info(ctx, fmt.Sprintf("  Description: %s", *sr.Description))
	} else {
		tflog.Info(ctx, "  Description: <nil>")
	}

	if sr.FirstPassword != nil {
		tflog.Info(ctx, fmt.Sprintf("  FirstPassword: [HIDDEN] (length: %d)", len(*sr.FirstPassword)))
	} else {
		tflog.Info(ctx, "  FirstPassword: <nil>")
	}

	if sr.CloudPanelID != nil {
		tflog.Info(ctx, fmt.Sprintf("  CloudPanelID: %s", *sr.CloudPanelID))
	} else {
		tflog.Info(ctx, "  CloudPanelID: <nil>")
	}

	// ✅ RSAKey (interface{})
	if sr.RSAKey != nil {
		tflog.Info(ctx, fmt.Sprintf("  RSAKey: %T = %v", sr.RSAKey, sr.RSAKey))
	} else {
		tflog.Info(ctx, "  RSAKey: <nil>")
	}

	// ✅ Objeto Datacenter
	tflog.Info(ctx, fmt.Sprintf("  Datacenter:"))
	tflog.Info(ctx, fmt.Sprintf("    ID: %s", sr.Datacenter.ID))
	tflog.Info(ctx, fmt.Sprintf("    CountryCode: %s", sr.Datacenter.CountryCode))
	tflog.Info(ctx, fmt.Sprintf("    Location: %s", sr.Datacenter.Location))

	// ✅ Objeto Status
	tflog.Info(ctx, fmt.Sprintf("  Status:"))
	tflog.Info(ctx, fmt.Sprintf("    State: %s", sr.Status.State))
	tflog.Info(ctx, fmt.Sprintf("    Percent: %d", sr.Status.Percent))

	// ✅ Objeto Image
	tflog.Info(ctx, fmt.Sprintf("  Image:"))
	tflog.Info(ctx, fmt.Sprintf("    ID: %s", sr.Image.ID))
	tflog.Info(ctx, fmt.Sprintf("    Name: %s", sr.Image.Name))

	// ✅ Objeto Hardware
	tflog.Info(ctx, fmt.Sprintf("  Hardware:"))
	tflog.Info(ctx, fmt.Sprintf("    VCore: %d", sr.Hardware.VCore))
	tflog.Info(ctx, fmt.Sprintf("    RAM: %d", sr.Hardware.RAM))
	tflog.Info(ctx, fmt.Sprintf("    CoresPerProcessor: %d", sr.Hardware.CoresPerProcessor))

	if sr.Hardware.BaremetalModelID != nil {
		tflog.Info(ctx, fmt.Sprintf("    BaremetalModelID: %s", *sr.Hardware.BaremetalModelID))
	} else {
		tflog.Info(ctx, "    BaremetalModelID: <nil>")
	}

	if sr.Hardware.FixedInstanceSizeID != nil {
		tflog.Info(ctx, fmt.Sprintf("    FixedInstanceSizeID: %s", *sr.Hardware.FixedInstanceSizeID))
	} else {
		tflog.Info(ctx, "    FixedInstanceSizeID: <nil>")
	}

	// ✅ HDDs en Hardware
	if len(sr.Hardware.HDDs) > 0 {
		tflog.Info(ctx, fmt.Sprintf("    HDDs: (count: %d)", len(sr.Hardware.HDDs)))
		for i, hdd := range sr.Hardware.HDDs {
			tflog.Info(ctx, fmt.Sprintf("      HDD[%d]:", i))
			tflog.Info(ctx, fmt.Sprintf("        ID: %s", hdd.ID))
			tflog.Info(ctx, fmt.Sprintf("        Size: %d", hdd.Size))
			tflog.Info(ctx, fmt.Sprintf("        IsMain: %t", hdd.IsMain))
			tflog.Info(ctx, fmt.Sprintf("        DiskType: %s", hdd.DiskType))
			tflog.Info(ctx, fmt.Sprintf("        DiskRaid: %s", hdd.DiskRaid))
		}
	} else {
		tflog.Info(ctx, "    HDDs: <empty>")
	}

	// ✅ Array IPs
	if len(sr.IPs) > 0 {
		tflog.Info(ctx, fmt.Sprintf("  IPs: (count: %d)", len(sr.IPs)))
		for i, ip := range sr.IPs {
			tflog.Info(ctx, fmt.Sprintf("    IP[%d]:", i))
			tflog.Info(ctx, fmt.Sprintf("      ID: %s", ip.ID))
			tflog.Info(ctx, fmt.Sprintf("      IP: %s", ip.IP))
			tflog.Info(ctx, fmt.Sprintf("      Type: %s", ip.Type))
			tflog.Info(ctx, fmt.Sprintf("      Main: %t", ip.Main))
			if ip.ReverseDNS != nil {
				tflog.Info(ctx, fmt.Sprintf("      ReverseDNS: %v", ip.ReverseDNS))
			} else {
				tflog.Info(ctx, "      ReverseDNS: <nil>")
			}
		}
	} else {
		tflog.Info(ctx, "  IPs: <empty>")
	}

	// ✅ Array PrivateNetworks
	if len(sr.PrivateNetworks) > 0 {
		tflog.Info(ctx, fmt.Sprintf("  PrivateNetworks: (count: %d)", len(sr.PrivateNetworks)))
		for i, pn := range sr.PrivateNetworks {
			tflog.Info(ctx, fmt.Sprintf("    PrivateNetwork[%d]:", i))
			tflog.Info(ctx, fmt.Sprintf("      ID: %s", pn.ID))
			tflog.Info(ctx, fmt.Sprintf("      Name: %s", pn.Name))
			tflog.Info(ctx, fmt.Sprintf("      ServerIP: %s", pn.ServerIP))
		}
	} else {
		tflog.Info(ctx, "  PrivateNetworks: <empty>")
	}

	// ✅ Objetos opcionales
	if sr.DVD != nil {
		tflog.Info(ctx, fmt.Sprintf("  DVD:"))
		tflog.Info(ctx, fmt.Sprintf("    ID: %s", sr.DVD.ID))
		tflog.Info(ctx, fmt.Sprintf("    Name: %s", sr.DVD.Name))
	} else {
		tflog.Info(ctx, "  DVD: <nil>")
	}

	if sr.Alerts != nil {
		tflog.Info(ctx, fmt.Sprintf("  Alerts: %+v", sr.Alerts))
	} else {
		tflog.Info(ctx, "  Alerts: <nil>")
	}

	if sr.MonitoringPolicy != nil {
		tflog.Info(ctx, fmt.Sprintf("  MonitoringPolicy:"))
		tflog.Info(ctx, fmt.Sprintf("    ID: %s", sr.MonitoringPolicy.ID))
		tflog.Info(ctx, fmt.Sprintf("    Name: %s", sr.MonitoringPolicy.Name))
	} else {
		tflog.Info(ctx, "  MonitoringPolicy: <nil>")
	}

	if sr.ConnectionSpeed != nil {
		tflog.Info(ctx, fmt.Sprintf("  ConnectionSpeed:"))
		if sr.ConnectionSpeed.Available != nil {
			tflog.Info(ctx, fmt.Sprintf("    Available: %v", sr.ConnectionSpeed.Available))
		} else {
			tflog.Info(ctx, "    Available: <nil>")
		}
		tflog.Info(ctx, fmt.Sprintf("    Current: %f", sr.ConnectionSpeed.Current))

		// Private connection speed
		if sr.ConnectionSpeed.Private != nil {
			tflog.Info(ctx, fmt.Sprintf("    Private:"))
			if sr.ConnectionSpeed.Private.Available != nil {
				tflog.Info(ctx, fmt.Sprintf("      Available: %v", sr.ConnectionSpeed.Private.Available))
			}
			tflog.Info(ctx, fmt.Sprintf("      Current: %f", sr.ConnectionSpeed.Private.Current))
		} else {
			tflog.Info(ctx, "    Private: <nil>")
		}

		// Public connection speed
		if sr.ConnectionSpeed.Public != nil {
			tflog.Info(ctx, fmt.Sprintf("    Public:"))
			if sr.ConnectionSpeed.Public.Available != nil {
				tflog.Info(ctx, fmt.Sprintf("      Available: %v", sr.ConnectionSpeed.Public.Available))
			}
			tflog.Info(ctx, fmt.Sprintf("      Current: %f", sr.ConnectionSpeed.Public.Current))
		} else {
			tflog.Info(ctx, "    Public: <nil>")
		}
	} else {
		tflog.Info(ctx, "  ConnectionSpeed: <nil>")
	}

	if sr.Redundancy != nil {
		tflog.Info(ctx, fmt.Sprintf("  Redundancy:"))
		tflog.Info(ctx, fmt.Sprintf("    Available: %t", sr.Redundancy.Available))
		tflog.Info(ctx, fmt.Sprintf("    Enabled: %t", sr.Redundancy.Enabled))
	} else {
		tflog.Info(ctx, "  Redundancy: <nil>")
	}

	if sr.Snapshot != nil {
		tflog.Info(ctx, fmt.Sprintf("  Snapshot:"))
		tflog.Info(ctx, fmt.Sprintf("    ID: %s", sr.Snapshot.ID))
		tflog.Info(ctx, fmt.Sprintf("    CreationDate: %s", sr.Snapshot.CreationDate))
		tflog.Info(ctx, fmt.Sprintf("    DeletionDate: %s", sr.Snapshot.DeletionDate))
	} else {
		tflog.Info(ctx, "  Snapshot: <nil>")
	}
}
