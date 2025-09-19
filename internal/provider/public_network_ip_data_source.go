package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicnetwork"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &PublicNetworkIpDataSource{}
	_ datasource.DataSourceWithConfigure = &PublicNetworkIpDataSource{}
)

func NewPublicNetworkIpDataSource() datasource.DataSource {
	return &PublicNetworkIpDataSource{}
}

type PublicNetworkIpDataSource struct {
	client *service.ApiPublicNetworkIpService
}

func (d *PublicNetworkIpDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network_ip"
}

func (d *PublicNetworkIpDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PublicNetworkIpDataSourceSchema(ctx)
}

func (d *PublicNetworkIpDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicNetworkIpService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicNetworkIpService, ok := client.(*service.ApiPublicNetworkIpService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = publicNetworkIpService
}

func (d *PublicNetworkIpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PublicNetworkIpModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	publicNetworkId := data.PublicNetworkId.ValueString()

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading IPs data source with ID %s in the public network %s", id, publicNetworkId))

	apiResponse, err := d.client.GetPublicNetworkIp(publicNetworkId, id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading IP in the public network",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving public network IP. Please report this issue to the provider developers.",
		)
		return
	}

	model, diags := models.NewPublicNetworkIpModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read IP with ID %s in the public network %s", id, publicNetworkId))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
