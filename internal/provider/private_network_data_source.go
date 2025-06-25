package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/privateNetwork"
)

var _ datasource.DataSource = &PrivateNetworkDataSource{}

func NewPrivateNetworkDataSource() datasource.DataSource {
	return &PrivateNetworkDataSource{}
}

type PrivateNetworkDataSource struct {
	client service.ApiPrivateNetworkService
}

func (d *PrivateNetworkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_network"
}

func (d *PrivateNetworkDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PrivateNetworkDataSourceSchema(ctx)
}

func (d *PrivateNetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = *privateNetworkService
}

func (d *PrivateNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PrivateNetworkModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Private Network Id",
			"Private network ID cannot be empty",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading private network with ID: %s", id))

	apiResponse, err := d.client.GetPrivateNetwork(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Private network with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading the private network",
			fmt.Sprintf("Could not read private network: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Private network with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	model, diags := models.NewPrivateNetwork(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
