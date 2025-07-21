package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/server"
)

var (
	_ datasource.DataSource              = &ServerDataSource{}
	_ datasource.DataSourceWithConfigure = &ServerDataSource{}
)

func NewServerDataSource() datasource.DataSource {
	return &ServerDataSource{}
}

type ServerDataSource struct {
	client *service.ApiServerService
}

func (d *ServerDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

func (d *ServerDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.ServerDataSourceSchema(ctx)
}

func (d *ServerDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetServerService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	serverService, ok := client.(*service.ApiServerService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = serverService
}

func (d *ServerDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.ServerDetailModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Server Id",
			"Server ID cannot be empty",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading server data source with ID: %s", id))

	apiResponse, err := d.client.GetServer(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server",
			fmt.Sprintf("Could not read server with ID %s: %s", id, err),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Server not found",
			fmt.Sprintf("Server with ID %s not found", id),
		)
		return
	}

	model, diags := models.NewServerDetailModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read server data source with ID: %s", id))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
