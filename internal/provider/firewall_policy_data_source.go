package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewallPolicy"
)

var (
	_ datasource.DataSource              = &FirewallPolicyDataSource{}
	_ datasource.DataSourceWithConfigure = &FirewallPolicyDataSource{}
)

func NewFirewallPolicyDataSource() datasource.DataSource {
	return &FirewallPolicyDataSource{}
}

type FirewallPolicyDataSource struct {
	client *service.ApiFirewallPolicyService
}

func (d *FirewallPolicyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_firewall_policy"
}

func (d *FirewallPolicyDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.FirewallPolicyDataSourceSchema(ctx)
}

func (d *FirewallPolicyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *FirewallPolicyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.FirewallPolicyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Firewall Policy Id",
			"Firewall Policy Id cannot be empty",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading Firewall Policy data source with ID: %s", id))

	apiResponse, err := d.client.GetFirewallPolicy(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading firewall policy",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Firewall policy not found",
			fmt.Sprintf("Firewall policy with ID %s not found", id),
		)
		return
	}

	model, diags := models.NewFirewallPolicyModel(ctx, *apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read firewall policy data source with ID: %s", id))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
