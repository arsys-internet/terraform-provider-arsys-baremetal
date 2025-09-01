package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/privateNetwork"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &PrivateNetworkServerDataSource{}

func NewPrivateNetworkServerDataSource() datasource.DataSource {
	return &PrivateNetworkServerDataSource{}
}

type PrivateNetworkServerDataSource struct {
	client service.ApiPrivateNetworkService
}

func (d *PrivateNetworkServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_network_server"
}

func (d *PrivateNetworkServerDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PrivateNetworkServerDataSourceSchema(ctx)
}

func (d *PrivateNetworkServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPrivateNetworkService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	privateNetworkService, ok := client.(*service.ApiPrivateNetworkService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *privateNetworkService
}

func (d *PrivateNetworkServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PrivateNetworkServerModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()
	privateNetworkId := data.PrivateNetworkId.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading private network server with Id: %s", id))

	apiResponse, err := d.client.GetPrivateNetworkServer(privateNetworkId, id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Private network Not Found",
				fmt.Sprintf("Private network server with id %s not found", id),
			)
			tflog.Info(ctx, fmt.Sprintf("Private network server with Id %s not found", id))
			return
		}

		resp.Diagnostics.AddError(
			"Error reading the private network server",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("Private network server not found"),
		)
		return
	}

	model, diags := models.NewPrivateNetworkServerFromResponse(ctx, privateNetworkId, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
