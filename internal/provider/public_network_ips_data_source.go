package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicnetwork"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &PublicNetworkIpsDataSource{}

func NewPublicNetworkIpsDataSource() datasource.DataSource {
	return &PublicNetworkIpsDataSource{}
}

type PublicNetworkIpsDataSource struct {
	client service.ApiPublicNetworkIpService
}

func (d *PublicNetworkIpsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network_ips"
}

func (d *PublicNetworkIpsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PublicNetworkIpsDataSourceSchema(ctx)
}

func (d *PublicNetworkIpsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicNetworkIpService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicNetworkService, ok := client.(*service.ApiPublicNetworkIpService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *publicNetworkService
}

func (d *PublicNetworkIpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PublicNetworkIpsModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetworkId := data.PublicNetworkId.ValueString()

	if publicNetworkId == "" {
		resp.Diagnostics.AddError(
			"Invalid public network Id",
			"Public network Id cannot be empty",
		)
		return
	}

	apiResponse, err := d.client.GetPublicNetworkIps(publicNetworkId)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading IPs in the public network",
			fmt.Sprintf("Could not read IPs in the public network: %s", err),
		)
		return
	}

	model, diags := models.NewPublicNetworkIps(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	model.PublicNetworkId = data.PublicNetworkId

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d IPs in the public network", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
