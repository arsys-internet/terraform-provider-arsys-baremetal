package provider

import (
	"os"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
)

// testAccProtoV6ProviderFactories is used to instantiate a provider during acceptance testing.
// The factory function is called for each Terraform CLI command to create a provider
// server that the CLI can connect to and interact with.
var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"cloudbuilder": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccProtoV6ProviderFactoriesWithEcho includes the echo provider alongside the scaffolding provider.
// It allows for testing assertions on data returned by an ephemeral resource during Open.
// The echoprovider is used to arrange tests by echoing ephemeral data into the Terraform state.
// This lets the data be referenced in test assertions with state checks.
var testAccProtoV6ProviderFactoriesWithEcho = map[string]func() (tfprotov6.ProviderServer, error){
	"cloudbuilder": providerserver.NewProtocol6WithError(New("test")()),
	"echo":         echoprovider.NewProviderServer(),
}

func testAccPreCheck(t *testing.T) {
	err := util.LoadEnv()

	if err != nil {
		t.Logf("Warning: %s", err)
		return
	}

	requiredEnvVars := []string{
		"CB_API_URL",
		"CB_API_KEY",
		"CB_USERNAME",
		"TF_ACC",
	}

	for _, env := range requiredEnvVars {
		if os.Getenv(env) == "" {
			t.Fatalf("The environment variable %s must be set", env)
		}
	}
}
