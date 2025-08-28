package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewallPolicy"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &FirewallPolicyRuleDataSource{}
	_ datasource.DataSourceWithConfigure = &FirewallPolicyRuleDataSource{}
)

func NewFirewallPolicyRuleDataSource() datasource.DataSource {
	return &FirewallPolicyRuleDataSource{}
}

type FirewallPolicyRuleDataSource struct {
	client *service.ApiFirewallPolicyService
}

func (d *FirewallPolicyRuleDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy_rule"
}

func (d *FirewallPolicyRuleDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyRuleDataSourceSchema(ctx)
}

func (d *FirewallPolicyRuleDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetFirewallPolicyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	policyService, ok := client.(*service.ApiFirewallPolicyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = policyService
}

func (d *FirewallPolicyRuleDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.FirewallPolicyRuleModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyId := data.FirewallPolicyId.ValueString()
	firewallPolicyRuleId := data.Id.ValueString()

	if firewallPolicyId == "" {
		resp.Diagnostics.AddError(
			"Invalid Firewall Policy Id",
			"firewall_policy_id cannot be empty",
		)
		return
	}

	if firewallPolicyRuleId == "" {
		resp.Diagnostics.AddError(
			"Invalid Firewall Policy Rule Id",
			"id cannot be empty",
		)
		return
	}

	apiResponse, err := d.client.GetFirewallPolicyRule(firewallPolicyId, firewallPolicyRuleId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading firewall policy rule",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Firewall policy not found",
			fmt.Sprintf("Firewall policy with Id %s not found", firewallPolicyId),
		)
		return
	}

	model := models.NewFirewallPolicyRuleModel(ctx, firewallPolicyId, *apiResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
