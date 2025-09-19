package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"

	service "terraform-provider-arsys-baremetal/internal/services/firewallpolicy"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &FirewallPolicyServerIPDataSource{}
	_ datasource.DataSourceWithConfigure = &FirewallPolicyServerIPDataSource{}
)

func NewFirewallPolicyServerIPDataSource() datasource.DataSource {
	return &FirewallPolicyServerIPDataSource{}
}

type FirewallPolicyServerIPDataSource struct {
	client *service.ApiFirewallPolicyService
}

func (d *FirewallPolicyServerIPDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy_server_ip"
}

func (d *FirewallPolicyServerIPDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyServerIpDataSourceSchema(ctx)
}

func (d *FirewallPolicyServerIPDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FirewallPolicyServerIPDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.FirewallPolicyServerIpModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	firewallPolicyId := data.FirewallPolicyId.ValueString()
	serverIpId := data.ServerIpId.ValueString()

	apiResponse, err := d.client.GetFirewallPolicyServerIP(firewallPolicyId, serverIpId)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Firewall Policy Not Found",
				fmt.Sprintf("Firewall policy server IP with ID %s not found in policy %s", serverIpId, firewallPolicyId),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading firewall policy",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving firewall policy server IP. Please report this issue to the provider developers.",
		)
		return
	}

	model := models.NewFirewallPolicyServerIpDataSourceModel(ctx, firewallPolicyId, *apiResponse)

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
