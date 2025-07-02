package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/publicIp"
)

var _ datasource.DataSource = &PublicIpDataSource{}

func NewPublicIpDataSource() datasource.DataSource {
	return &PublicIpDataSource{}
}

type PublicIpDataSource struct {
	client service.ApiPublicIpService
}

func (d *PublicIpDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (d *PublicIpDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.PublicIpDataSourceSchema(ctx)
}

func (d *PublicIpDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetPublicIpService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	publicIpService, ok := client.(*service.ApiPublicIpService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *publicIpService
}

func (d *PublicIpDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.PublicIpModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid public IP Id",
			"Public IP ID cannot be empty",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading public IP with ID: %s", id))

	apiResponse, err := d.client.GetPublicIp(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Public IP with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading the public IP",
			fmt.Sprintf("Could not read public IP: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Public IP with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	model, diags := models.NewPublicIp(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
