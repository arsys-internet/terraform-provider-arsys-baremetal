package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/sshkey"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &SshKeyDataSource{}
	_ datasource.DataSourceWithConfigure = &SshKeyDataSource{}
)

func NewSshKeyDataSource() datasource.DataSource {
	return &SshKeyDataSource{}
}

type SshKeyDataSource struct {
	client *service.ApiSshKeyService
}

func (d *SshKeyDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_key"
}

func (d *SshKeyDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.SshKeyDataSourceSchema(ctx)
}

func (d *SshKeyDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetSshKeyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	sshKeyService, ok := client.(*service.ApiSshKeyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			fmt.Sprintf("An internal error occurred. Please report this issue to the provider developers."),
		)
		return
	}

	d.client = sshKeyService
}

func (d *SshKeyDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data models.SshKeyModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := data.Id.ValueString()

	tflog.Info(ctx, fmt.Sprintf("Reading SSH key data source with ID: %s", id))

	apiResponse, err := d.client.GetSshKey(id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SSH key",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	if apiResponse == nil {
		resp.Diagnostics.AddError(
			"Internal Error",
			"An unexpected error occurred while retrieving SSH key. Please try again or report this issue to the provider developers",
		)
		return
	}

	model, diags := models.NewSshKeyModel(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read SSH key data source with ID: %s", id))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
