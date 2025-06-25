package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"os"
	"terraform-provider-arsys-baremetal/internal/client"
	"testing"
)

// TestAccProtoV6ProviderFactories are used to instantiate a provider during acceptance testing.
// Exported for use in other testing packages
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"arsys-baremetal": providerserver.NewProtocol6WithError(New("test")()),
}

// testAccProtoV6ProviderFactories is the local version for use in this package
var testAccProtoV6ProviderFactories = TestAccProtoV6ProviderFactories

// TestAccPreCheck verifies that required environment variables are configured
// Exported for use in other testing packages
func TestAccPreCheck(t *testing.T) {
	// Check required environment variables
	if v := os.Getenv("BAREMETAL_HOST"); v == "" {
		t.Fatal("BAREMETAL_HOST must be set for acceptance tests")
	}
	if v := os.Getenv("BAREMETAL_API_TOKEN"); v == "" {
		t.Fatal("BAREMETAL_API_TOKEN must be set for acceptance tests")
	}
}

func getTestClient() *client.APIClient {
	return client.NewAPIClient(
		os.Getenv("BAREMETAL_API_TOKEN"),
		os.Getenv("BAREMETAL_HOST"),
	)
}
