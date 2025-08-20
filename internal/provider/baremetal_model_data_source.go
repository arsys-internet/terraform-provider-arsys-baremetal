package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/server"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &BaremetalModelDataSource{}

func NewBaremetalModelDataSource() datasource.DataSource {
	return &BaremetalModelDataSource{}
}

type BaremetalModelDataSource struct {
	client service.ApiServerService
}

func (d *BaremetalModelDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_baremetal_model"
}

func (d *BaremetalModelDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.BaremetalModelDataSourceSchema(ctx)
}

func (d *BaremetalModelDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetServerService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	baremetalService, ok := client.(*service.ApiServerService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *baremetalService
}

func (d *BaremetalModelDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.BaremetalModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Baremetal Model Id",
			"Baremetal model ID cannot be empty",
		)
		return
	}

	apiResponse, err := d.client.GetBaremetalModel(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Baremetal Model Not Found",
				fmt.Sprintf("Error: %s", err.Error()),
			)
			tflog.Info(ctx, fmt.Sprintf("Error: %s", err.Error()))
			return
		}

		resp.Diagnostics.AddError(
			"Error reading baremetal model",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Not Found",
			fmt.Sprintf("Baremetal model not found"),
		)
		return
	}

	model, diags := models.NewBaremetalModel(ctx, *apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
