package provider

import (
	"fmt"
	"regexp"
	"strings"
	service "terraform-provider-arsys-baremetal/internal/services/privatenetwork"
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

func TestAccPrivateNetworkResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
			time.Sleep(5 * time.Second)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		CheckDestroy: testAccCheckPrivateNetworkDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNetworkResourceConfig(
					"acc-test-pn-terraform",
					"81DEF28500FBC2A973FC0C620DF5B721",
					"10.0.1.0",
					"255.255.255.0",
					"Private network for acceptance testing",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("acc-test-pn-terraform"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Private network for acceptance testing"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("network_address"),
						knownvalue.StringExact("10.0.1.0"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("subnet_mask"),
						knownvalue.StringExact("255.255.255.0"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("cloudpanel_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter").AtMapKey("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("servers"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
			{
				ResourceName:      "arsys-baremetal_private_network.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"description",
				},
			},
			{
				Config: testAccPrivateNetworkResourceConfig(
					"acc-test-pn-terraform-updated",
					"81DEF28500FBC2A973FC0C620DF5B721",
					"10.0.2.0",
					"255.255.255.0",
					"Updated private network description",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("acc-test-pn-terraform-updated"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated private network description"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("network_address"),
						knownvalue.StringExact("10.0.2.0"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("subnet_mask"),
						knownvalue.StringExact("255.255.255.0"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccCheckPrivateNetworkDestroy(s *terraform.State) error {
	testService := getTestPrivateNetworkService()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "arsys-baremetal_private_network" {
			continue
		}

		id := rs.Primary.ID
		if id == "" {
			continue
		}

		_, err := testService.GetPrivateNetwork(id)

		if err == nil {
			return fmt.Errorf("private network %s still exists", id)
		}

		if strings.Contains(err.Error(), "not found") {
			continue
		}

		return fmt.Errorf("error checking private network %s: %s", id, err)
	}

	return nil
}

func testAccPrivateNetworkResourceConfig(name, datacenterId, networkAddress, subnetMask, description string) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_private_network" "test" {
  name            = %[1]q
  datacenter_id   = %[2]q
  network_address = %[3]q
  subnet_mask     = %[4]q
  description     = %[5]q
}
`, name, datacenterId, networkAddress, subnetMask, description)
}

func getTestPrivateNetworkService() *service.ApiPrivateNetworkService {
	apiClient := getTestClient()
	return service.NewApiPrivateNetworkService(apiClient)
}
