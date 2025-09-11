package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/privateNetwork"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &PrivateNetworkServersDataSource{}

func NewPrivateNetworkServersDataSource() datasource.DataSource {
	return &PrivateNetworkServersDataSource{}
}

type PrivateNetworkServersDataSource struct {
	client *service.ApiPrivateNetworkService
}

func (d *PrivateNetworkServersDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_private_network_servers"
}

func (d *PrivateNetworkServersDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PrivateNetworkServersDataSourceSchema(ctx)
}

func (d *PrivateNetworkServersDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

	d.client = privateNetworkService
}

func (d *PrivateNetworkServersDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading all private network servers")
	var data models.PrivateNetworkServersModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	apiResponse, err := d.client.GetPrivateNetworkServers(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading private network servers",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	model, diags := models.NewPrivateNetworkServers(ctx, id, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d private network servers", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
