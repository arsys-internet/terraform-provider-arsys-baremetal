package provider

import (
	"context"
	"errors"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicip"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

var _ datasource.DataSource = &SubnetDataSource{}

func NewSubnetDataSource() datasource.DataSource {
	return &SubnetDataSource{}
}

type SubnetDataSource struct {
	client service.ApiPublicIpService
}

func (d *SubnetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_subnet"
}

func (d *SubnetDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.SubnetDataSourceSchema(ctx)
}

func (d *SubnetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *SubnetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PublicIpModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	apiResponse, err := d.client.GetSubnet(id)
	if err != nil {
		if errors.Is(err, util.ErrNotFound) {
			resp.Diagnostics.AddError(
				"Subnet not found",
				fmt.Sprintf("Subnet with ID %s was not found", id),
			)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading subnet",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving subnet. Please try again or report this issue to the provider developers",
		)
		return
	}

	model, diags := models.NewPublicIpModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
