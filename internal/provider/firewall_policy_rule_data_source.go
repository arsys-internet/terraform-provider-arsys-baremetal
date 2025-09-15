package provider

import (
	"context"
	"fmt"
	"strings"
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

	apiResponse, err := d.client.GetFirewallPolicyRule(firewallPolicyId, firewallPolicyRuleId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Firewall Policy Rule Not Found",
				fmt.Sprintf("Firewall policy rule with ID %s not found in policy %s", firewallPolicyRuleId, firewallPolicyId),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading firewall policy rule",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving firewall policy rule. Please report this issue to the provider developers.",
		)
		return
	}

	model := models.NewFirewallPolicyRuleModel(ctx, firewallPolicyId, *apiResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
