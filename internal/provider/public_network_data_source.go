package provider

import (
	"context"
	"fmt"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicNetwork"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &PublicNetworkDataSource{}

func NewPublicNetworkDataSource() datasource.DataSource {
	return &PublicNetworkDataSource{}
}

type PublicNetworkDataSource struct {
	client service.ApiPublicNetworkService
}

func (d *PublicNetworkDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_network"
}

func (d *PublicNetworkDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PublicNetworkDataSourceSchema(ctx)
}

func (d *PublicNetworkDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *PublicNetworkDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PublicNetworkModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid public network Id",
			"Public network ID cannot be empty",
		)
		return
	}

	apiResponse, err := d.client.GetPublicNetwork(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			resp.Diagnostics.AddError(
				"Public network not found",
				fmt.Sprintf("Public network with ID %s was not found", id),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading the public network",
			fmt.Sprintf("Could not read public network: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Public network not found",
			fmt.Sprintf("Public network with ID %s was not found", id),
		)
		return
	}

	model, diags := models.NewPublicNetworkModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
