package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	firewallpolicy "terraform-provider-arsys-baremetal/internal/models/firewall_policy"

	service "terraform-provider-arsys-baremetal/internal/services/firewall_policy"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &FirewallPolicyResource{}
	_ resource.ResourceWithConfigure   = &FirewallPolicyResource{}
	_ resource.ResourceWithImportState = &FirewallPolicyResource{}
)

func NewFirewallPolicyResource() resource.Resource {
	return &FirewallPolicyResource{}
}

type FirewallPolicyResource struct {
	client *service.ApiFirewallPolicyService
}

func (r *FirewallPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy"
}

func (r *FirewallPolicyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyResourceSchema(ctx)
}

func (r *FirewallPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetFirewallPolicyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	policyService, ok := client.(*service.ApiFirewallPolicyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	r.client = policyService
}

func (r *FirewallPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data models.FirewallPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if data.Name.IsNull() || data.Name.ValueString() == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("name"),
			"Missing required field",
			"'name' field is required when creating a firewall policy",
		)
	}

	if data.Rules.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("rules"),
			"Missing required field",
			"'rules' field is required when creating a firewall policy",
		)
	} else {
		rulesDiags := firewallpolicy.ValidateFirewallRules(data.Rules, path.Root("rules"))
		resp.Diagnostics.Append(rulesDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createRequest, err := data.ToCreateRequest()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating firewall policy",
			fmt.Sprintf("Could not create firewall policy: %s", err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Creating firewall policy: %s", createRequest.Name))

	apiResponse, err := r.client.CreateFirewallPolicy(&createRequest)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("FIREWALL_POLICY")

	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
		[]string{util.StatusConfiguring},
		[]string{util.StateActive},
	)

	_, diags := util.WaitForResourceState(
		ctx,
		apiResponse.Id,
		r.client,
		waitOptions,
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Wait for firewall policy state failed")
		return
	}

	finalFirewallPolicy, err := r.client.GetFirewallPolicy(apiResponse.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting final firewall policy state",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	finalModel, diags := models.NewFirewallPolicyModel(ctx, *finalFirewallPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		tflog.Error(ctx, "Failed to create final resource model")
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Created firewall policy with Id: %s", finalModel.Id.ValueString()))

	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *FirewallPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data models.FirewallPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading firewall policy with Id: %s", id))

	apiResponse, err := r.client.GetFirewallPolicy(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Firewall policy with Id %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Firewall policy with Id %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	readModel, diags := models.NewFirewallPolicyModelFromRead(ctx, apiResponse, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, readModel)...)
}

func (r *FirewallPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.FirewallPolicyModel
	var state models.FirewallPolicyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.Id.ValueString()
	tflog.Info(ctx, fmt.Sprintf("Updating firewall policy with Id: %s", id))

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

	_, err := r.client.UpdateFirewallPolicy(id, &updateRequest)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	finalModel := state

	if !plan.Name.Equal(state.Name) {
		finalModel.Name = plan.Name
	}

	if !plan.Description.Equal(state.Description) {
		finalModel.Description = plan.Description
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully updated firewall policy with Id: %s", id))
	resp.Diagnostics.Append(resp.State.Set(ctx, finalModel)...)
}

func (r *FirewallPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data models.FirewallPolicyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Deleting firewall policy with Id: %s", id))

	err := r.client.DeleteFirewallPolicy(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Firewall Policy %s was already deleted", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error deleting firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	timeouts := util.GetResourceTimeouts("FIREWALL_POLICY")

	waitOptions := util.NewWaitOptions(
		timeouts.Default,
		timeouts.RetryInterval,
		timeouts.MinTimeout,
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

	tflog.Info(ctx, fmt.Sprintf("Deleted firewall policy with Id: %s", id))
}

func (r *FirewallPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
