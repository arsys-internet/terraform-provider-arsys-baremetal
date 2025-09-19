package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/serverappliance"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &ServerAppliancesDataSource{}

func NewServerAppliancesDataSource() datasource.DataSource {
	return &ServerAppliancesDataSource{}
}

type ServerAppliancesDataSource struct {
	client *service.ApiServerApplianceService
}

func (d *ServerAppliancesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_appliances"
}

func (d *ServerAppliancesDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.ServerAppliancesDataSourceSchema(ctx)
}

func (d *ServerAppliancesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetServerApplianceService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	serverApplianceService, ok := client.(*service.ApiServerApplianceService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = serverApplianceService
}

func (d *ServerAppliancesDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading all server appliances")

	apiResponse, err := d.client.GetServerAppliances()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading server appliances",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	model, diags := models.NewServerAppliances(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d server appliances", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
