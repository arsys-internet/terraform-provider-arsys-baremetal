package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/privateNetwork"
)

var _ datasource.DataSource = &PrivateNetworksDataSource{}

func NewPrivateNetworksDataSource() datasource.DataSource {
	return &PrivateNetworksDataSource{}
}

type PrivateNetworksDataSource struct {
	client *service.ApiPrivateNetworkService
}

func (d *PrivateNetworksDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_networks"
}

func (d *PrivateNetworksDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PrivateNetworksDataSourceSchema(ctx)
}

func (d *PrivateNetworksDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPrivateNetworkService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	privateNetworkService, ok := client.(*service.ApiPrivateNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *pNetwork.ApiPrivateNetworkService, got: %T. Please report this issue to the provider developers.", client),
		)
		return
	}

	d.client = privateNetworkService
}

func (d *PrivateNetworksDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading all private networks")

	apiResponse, err := d.client.GetPrivateNetworks()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private networks",
			fmt.Sprintf("Could not read private networks: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, "No private networks found")
		apiResponse = []models.PrivateNetworkResponse{}
	}

	model, diags := models.NewPrivateNetworks(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d private networks", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
