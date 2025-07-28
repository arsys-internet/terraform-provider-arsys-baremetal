package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/server"
)

var _ datasource.DataSource = &BaremetalModelsDataSource{}
var _ datasource.DataSourceWithConfigure = &BaremetalModelsDataSource{}

func NewBaremetalModelsDataSource() datasource.DataSource {
	return &BaremetalModelsDataSource{}
}

type BaremetalModelsDataSource struct {
	client *service.ApiServerService
}

func (b *BaremetalModelsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_models"
}

func (b *BaremetalModelsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.BaremetalModelsDataSourceSchema(ctx)
}

func (b *BaremetalModelsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetServerService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"An internal error occurred. Please report this issue to the provider developers.",
		)
		return
	}

	serverService, ok := client.(*service.ApiServerService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"An internal error occurred. Please report this issue to the provider developers.",
		)
		return
	}

	b.client = serverService
}

func (b *BaremetalModelsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading all baremetal models")

	apiResponse, err := b.client.GetBaremetalModels()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading baremetal models",
			fmt.Sprintf("Error: %s", err),
		)
		return
	}

	model, diags := models.NewBaremetalModels(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d baremetal models", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
