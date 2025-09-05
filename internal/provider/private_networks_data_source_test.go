package provider

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccPrivateNetworksDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
			time.Sleep(5 * time.Second)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPrivateNetworksResourceConfig(),
			},
			{
				Config: testAccPrivateNetworksDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_private_networks.all", "id"),
					testAccCheckPrivateNetworkInList("pn-list-test-terraform"),
				),
			},
		},
	})
}

func testAccCheckPrivateNetworkInList(targetName string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_private_networks.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["private_networks.#"]
		if !ok {
			return fmt.Errorf("private_networks count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid private_networks count '%s': %v", countStr, err)
		}

		if count == 0 {
			return fmt.Errorf("expected at least 1 private network, got 0")
		}

		for i := 0; i < count; i++ {
			nameKey := fmt.Sprintf("private_networks.%d.name", i)
			name := rs.Primary.Attributes[nameKey]

			if name == targetName {
				return validatePrivateNetworkEssentials(rs.Primary.Attributes, i, targetName)
			}
		}

		return fmt.Errorf("private network '%s' not found in list of %d elements", targetName, count)
	}
}

func validatePrivateNetworkEssentials(attributes map[string]string, index int, targetName string) error {
	idKey := fmt.Sprintf("private_networks.%d.id", index)
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("private network '%s' missing ID", targetName)
	}

	nameKey := fmt.Sprintf("private_networks.%d.name", index)
	if name, exists := attributes[nameKey]; !exists {
		return fmt.Errorf("private network '%s' missing name", targetName)
	} else if name != targetName {
		return fmt.Errorf("expected name '%s', got '%s'", targetName, name)
	}

	stateKey := fmt.Sprintf("private_networks.%d.state", index)
	if state, exists := attributes[stateKey]; !exists {
		return fmt.Errorf("private network '%s' missing state", targetName)
	} else if state != "ACTIVE" {
		return fmt.Errorf("private network '%s' expected state ACTIVE, got '%s'", targetName, state)
	}

	serversKey := fmt.Sprintf("private_networks.%d.servers.#", index)
	if _, exists := attributes[serversKey]; !exists {
		return fmt.Errorf("private network '%s' missing servers list", targetName)
	}

	networkAddressKey := fmt.Sprintf("private_networks.%d.network_address", index)
	if networkAddr, exists := attributes[networkAddressKey]; !exists || networkAddr == "" {
		return fmt.Errorf("private network '%s' missing network_address", targetName)
	}

	return nil
}

func testAccPrivateNetworksResourceConfig() string {
	return `
resource "arsys-baremetal_private_network" "test" {
  name            = "pn-list-test-terraform"
  datacenter_id   = "81DEF28500FBC2A973FC0C620DF5B721"
  network_address = "10.0.3.0"
  subnet_mask     = "255.255.255.0"
  description     = "Private network for list testing"
}
`
}

func testAccPrivateNetworksDataSourceConfig() string {
	return testAccPrivateNetworksResourceConfig() + `
data "arsys-baremetal_private_networks" "all" {
}
`
}
