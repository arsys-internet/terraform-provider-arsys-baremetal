package provider

import (
	"os"
	"terraform-provider-arsys-baremetal/internal/client"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"arsys-baremetal": providerserver.NewProtocol6WithError(New("test")()),
}

var testAccProtoV6ProviderFactories = TestAccProtoV6ProviderFactories

func TestAccPreCheck(t *testing.T) {
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
