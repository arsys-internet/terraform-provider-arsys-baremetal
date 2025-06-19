package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"os"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/util"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure BaremetalProvider satisfies various provider interfaces.
var _ provider.Provider = &BaremetalProvider{}

// BaremetalProvider defines the provider implementation.
type BaremetalProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// BaremetalProviderModel describes the provider data model.
type BaremetalProviderModel struct {
	Host  types.String `tfsdk:"host"`
	Token types.String `tfsdk:"token"`
}

func (p *BaremetalProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "arsys-baremetal"
	resp.Version = p.version
}

func (p *BaremetalProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Baremetal provider.",
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional:    true,
				Description: "Host of the Baremetal API.",
			},
			"token": schema.StringAttribute{
				Optional:    true,
				Description: "API key of the Baremetal API.",
			},
		},
	}
}

func (p *BaremetalProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Baremetal provider client")

	// Retrieve provider data from configuration
	var config BaremetalProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Unknown Baremetal API Host",
			"The provider cannot create the Baremetal API client as there is an unknown configuration value for the Baremetal API host. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAREMETAL_HOST environment variable.",
		)
	}

	if config.Token.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Unknown API Token",
			"The provider cannot create the Baremetal API client as there is an unknown configuration value for the  API token. "+
				"Either target apply the source of the value first, set the value statically in the configuration, or use the BAREMETAL_API_TOKEN environment variable.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	var (
		host     string
		apiToken string
	)

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	} else {
		host = os.Getenv("BAREMETAL_HOST")
	}

	if !config.Token.IsNull() {
		apiToken = config.Token.ValueString()
	} else {
		apiToken = os.Getenv("BAREMETAL_API_TOKEN")
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"),
			"Missing Baremetal API Host",
			"The provider cannot create the Baremetal API client as there is a missing or empty value for the Baremetal API host. "+
				"Set the host value in the configuration or use the BAREMETAL_HOST environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if apiToken == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("token"),
			"Missing Baremetal API Token",
			"The provider cannot create the Baremetal API client as there is a missing or empty value for the Baremetal API token. "+
				"Set the token value in the configuration or use the BAREMETAL_API_TOKEN environment variable. "+
				"If either is already set, ensure the value is not empty.",
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating client")

	// Create a new client using the configuration values
	c := client.NewAPIClient(apiToken, host)

	// Make the client available during DataSource and Resource type Configure methods.
	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured client", map[string]any{"success": true})
}

func (p *BaremetalProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExecutionGroupResource,
	}
}

func (p *BaremetalProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *BaremetalProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewPrivateNetworkDataSource,
	}
}

func (p *BaremetalProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &BaremetalProvider{
			version: version,
		}
	}
}

func init() {
	if err := util.LoadEnv(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}
