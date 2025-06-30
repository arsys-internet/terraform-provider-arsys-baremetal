package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strings"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/datacenter"
)

var _ datasource.DataSource = &DatacenterDataSource{}

func NewDatacenterDataSource() datasource.DataSource {
	return &DatacenterDataSource{}
}

type DatacenterDataSource struct {
	client service.ApiDatacenterService
}

func (d *DatacenterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_datacenter"
}

func (d *DatacenterDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.DatacenterDataSourceSchema(ctx)
}

func (d *DatacenterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetDatacenterService(req.ProviderData)

	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	privateNetworkService, ok := client.(*service.ApiDatacenterService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = *privateNetworkService
}

func (d *DatacenterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.DatacenterModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.ID.ValueString()

	if id == "" {
		resp.Diagnostics.AddError(
			"Invalid Datacenter Id",
			"datacenter_id cannot be empty",
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Reading datacenter with ID: %s", id))

	apiResponse, err := d.client.GetDatacenter(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			tflog.Info(ctx, fmt.Sprintf("Datacenter with ID %s not found, removing from state", id))
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error reading the private network",
			fmt.Sprintf("Could not read private network: %s", err),
		)
		return
	}

	if apiResponse == nil {
		tflog.Info(ctx, fmt.Sprintf("Datacenter with ID %s not found, removing from state", id))
		resp.State.RemoveResource(ctx)
		return
	}

	model, diags := models.NewDatacenter(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
