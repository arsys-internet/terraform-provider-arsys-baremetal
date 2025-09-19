package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/firewallpolicy"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading Firewall Policy data source with ID: %s", id))

	apiResponse, err := d.client.GetFirewallPolicy(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Firewall Policy Not Found",
				fmt.Sprintf("Firewall Policy with ID %s not found", id),
			)
			tflog.Info(ctx, fmt.Sprintf("Firewall Policy with ID %s not found", id))
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
			"An unexpected error occurred while retrieving firewall policy. Please try again or report this issue to the provider developers",
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
