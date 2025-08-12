package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicIp"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &PublicIpsDataSource{}

func NewPublicIpsDataSource() datasource.DataSource {
	return &PublicIpsDataSource{}
}

type PublicIpsDataSource struct {
	client service.ApiPublicIpService
}

func (d *PublicIpsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ips"
}

func (d *PublicIpsDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PublicIpsDataSourceSchema(ctx)
}

func (d *PublicIpsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicIpService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicIpService, ok := client.(*service.ApiPublicIpService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Datasource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *publicIpService
}

func (d *PublicIpsDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	apiResponse, err := d.client.GetPublicIps()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading public IPs",
			fmt.Sprintf("Could not read public IPs: %s", err),
		)
		return
	}

	model, diags := models.NewPublicIps(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d public IPs", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
