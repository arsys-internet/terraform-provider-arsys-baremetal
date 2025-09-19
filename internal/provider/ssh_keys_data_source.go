package provider

import (
	"context"
	"fmt"
	"terraform-provider-arsys-baremetal/internal/models"
	service "terraform-provider-arsys-baremetal/internal/services/sshkey"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ datasource.DataSource = &SshKeysDataSource{}
var _ datasource.DataSourceWithConfigure = &SshKeysDataSource{}

func NewSshKeysDataSource() datasource.DataSource {
	return &SshKeysDataSource{}
}

type SshKeysDataSource struct {
	client *service.ApiSshKeyService
}

func (d *SshKeysDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ssh_keys"
}

func (d *SshKeysDataSource) Schema(ctx context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = models.SshKeysDataSourceSchema(ctx)
}

func (d *SshKeysDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client := service.GetSshKeyService(req.ProviderData)
	if client == nil {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"An internal error occurred. Please report this issue to the provider developers.",
		)
		return
	}

	policyService, ok := client.(*service.ApiSshKeyService)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected DataSource Configure Type",
			"An internal error occurred. Please report this issue to the provider developers.",
		)
		return
	}

	d.client = policyService
}

func (d *SshKeysDataSource) Read(ctx context.Context, _ datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Reading all SSH keys")

	apiResponse, err := d.client.GetSshKeys()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SSH keys",
			fmt.Sprintf("Error: %s", err.Error()),
		)
		return
	}

	model, diags := models.NewSshKeys(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Successfully read %d SSH keys", len(apiResponse)))

	resp.Diagnostics.Append(resp.State.Set(ctx, model)...)
}
