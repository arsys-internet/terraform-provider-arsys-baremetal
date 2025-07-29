package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewallPolicy"
)

var _ datasource.DataSource = &FirewallPoliciesDataSource{}
var _ datasource.DataSourceWithConfigure = &FirewallPoliciesDataSource{}

func NewFirewallPoliciesDataSource() datasource.DataSource {
	return &FirewallPoliciesDataSource{}
}

type FirewallPoliciesDataSource struct {
	client *service.ApiFirewallPolicyService
}

func (d *FirewallPoliciesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policies"
}

func (d *FirewallPoliciesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.FirewallPoliciesDataSourceSchema(ctx)
}

func (d *FirewallPoliciesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetFirewallPolicyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"An internal error occurred. Please report this issue to the provider developers.",
		)
		return
	}

	policyService, ok := client.(*service.ApiFirewallPolicyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"An internal error occurred. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = policyService
}

func (d *FirewallPoliciesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading all firewall policies")

	apiResponse, err := d.client.GetFirewallPolicies()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading firewall policies",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	model, diags := models.NewFirewallPolicies(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d firewall policies", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
