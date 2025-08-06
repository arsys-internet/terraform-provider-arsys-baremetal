package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"log"
	"os"
	"terraform-provider-arsys-baremetal/internal/client"
	"terraform-provider-arsys-baremetal/internal/util"
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

func (p *BaremetalProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "arsys-baremetal"
	resp.Version = p.version
}

func (p *BaremetalProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
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
		if host == "" {
			host = "https://api.cloudbuilder.es/v1"
		}
	}

	if !config.Token.IsNull() {
		apiToken = config.Token.ValueString()
	} else {
		apiToken = os.Getenv("BAREMETAL_API_TOKEN")
	}

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

	c := client.NewAPIClient(apiToken, host)

	resp.DataSourceData = c
	resp.ResourceData = c

	tflog.Info(ctx, "Configured client", map[string]any{"success": true})
}

func (p *BaremetalProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewFirewallPolicyResource,
		NewPrivateNetworkResource,
		NewPublicIpResource,
		NewPublicNetworkResource,
		NewPublicNetworkServerResource,
	}
}

func (p *BaremetalProvider) EphemeralResources(_ context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *BaremetalProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewDatacenterDataSource,
		NewDatacentersDataSource,
		NewFirewallPolicyDataSource,
		NewFirewallPoliciesDataSource,
		NewFirewallPolicyServerIPDataSource,
		NewFirewallPolicyServerIPsDataSource,
		NewPrivateNetworkDataSource,
		NewPrivateNetworksDataSource,
		NewPublicIpDataSource,
		NewPublicIpsDataSource,
		NewPublicNetworkDataSource,
		NewPublicNetworksDataSource,
		NewServerApplianceDataSource,
		NewServerAppliancesDataSource,
	}
}

func (p *BaremetalProvider) Functions(_ context.Context) []func() function.Function {
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
