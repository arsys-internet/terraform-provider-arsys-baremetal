package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicnetwork"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &PublicNetworksDataSource{}

func NewPublicNetworksDataSource() datasource.DataSource {
	return &PublicNetworksDataSource{}
}

type PublicNetworksDataSource struct {
	client service.ApiPublicNetworkService
}

func (d *PublicNetworksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_networks"
}

func (d *PublicNetworksDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PublicNetworksDataSourceSchema(ctx)
}

func (d *PublicNetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicNetworkService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicNetworkService, ok := client.(*service.ApiPublicNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *publicNetworkService
}

func (d *PublicNetworksDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	apiResponse, err := d.client.GetPublicNetworks()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading public networks",
			fmt.Sprintf("Could not read public networks: %s", err),
		)
		return
	}

	model, diags := models.NewPublicNetworks(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d public networks", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
