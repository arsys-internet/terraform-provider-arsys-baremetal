package provider

import (
	"fmt"
	"regexp"
	"strings"
	service "terraform-provider-arsys-baremetal/internal/services/publicnetwork"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccPublicNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
			time.Sleep(5 * time.Second)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		CheckDestroy: testAccCheckPublicNetworkDestroy,
		Steps: []resource.TestStep{
			// Test Create
			{
				Config: testAccPublicNetworkResourceConfig(
					"test-public-network",
					"Test public network for acceptance testing",
					"81DEF28500FBC2A973FC0C620DF5B721",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("public_name"),
						knownvalue.StringExact("test-public-network"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Test public network for acceptance testing"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("start_date"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("same_vlan"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("servers"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("ips"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("availability_zones"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("last_logs"),
						knownvalue.NotNull(),
					),
				},
			},
			// Test Import
			{
				ResourceName:      "arsys-baremetal_public_network.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			// Test Update
			{
				Config: testAccPublicNetworkResourceConfig(
					"updated-public-network",
					"Updated public network description",
					"81DEF28500FBC2A973FC0C620DF5B721",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("public_name"),
						knownvalue.StringExact("updated-public-network"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated public network description"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_network.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
				},
			},
		},
	})
}

func testAccCheckPublicNetworkDestroy(s *terraform.State) error {
	testService := getTestPublicNetworkService()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "arsys-baremetal_public_network" {
			continue
		}

		id := rs.Primary.ID
		if id == "" {
			continue
		}

		_, err := testService.GetPublicNetwork(id)

		if err == nil {
			return fmt.Errorf("public network %s still exists", id)
		}

		if strings.Contains(err.Error(), "not found") {
			continue
		}

		return fmt.Errorf("error checking public network %s: %s", id, err)
	}

	return nil
}

func testAccPublicNetworkResourceConfig(publicName, description, datacenterID string) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_public_network" "test" {
  public_name   = %[1]q
  description   = %[2]q
  datacenter_id = %[3]q
}
`, publicName, description, datacenterID)
}

func getTestPublicNetworkService() *service.ApiPublicNetworkService {
	apiClient := getTestClient()
	return service.NewApiPublicNetworkService(apiClient)
}
