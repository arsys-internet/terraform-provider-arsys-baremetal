package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewallpolicy"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var (
	_ datasource.DataSource              = &FirewallPolicyServerIPsDataSource{}
	_ datasource.DataSourceWithConfigure = &FirewallPolicyServerIPsDataSource{}
)

func NewFirewallPolicyServerIPsDataSource() datasource.DataSource {
	return &FirewallPolicyServerIPsDataSource{}
}

type FirewallPolicyServerIPsDataSource struct {
	client *service.ApiFirewallPolicyService
}

func (d *FirewallPolicyServerIPsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy_server_ips"
}

func (d *FirewallPolicyServerIPsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyServerIPsSchema(ctx)
}

func (d *FirewallPolicyServerIPsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FirewallPolicyServerIPsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.FirewallPolicyServerIpsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	apiResponse, err := d.client.GetFirewallPolicy(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading firewall policy server IP",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving firewall policy server IPs. Please report this issue to the provider developers.",
		)
		return
	}

	model, diags := models.NewFirewallPolicyServerIpsModel(ctx, id, apiResponse.ServerIPs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
