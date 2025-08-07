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

func TestAccPublicNetworksDataSource(t *testing.T) {
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
				Config: testAccPublicNetworksResourceConfig(),
			},
			{
				Config: testAccPublicNetworksDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_public_networks.all", "id"),
					testAccCheckPublicNetworksListNotEmpty(),
				),
			},
		},
	})
}

func testAccCheckPublicNetworksListNotEmpty() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_public_networks.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["public_networks.#"]
		if !ok {
			return fmt.Errorf("public_networks count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid public_networks count '%s': %v", countStr, err)
		}

		if count == 0 {
			return fmt.Errorf("expected at least 1 public network, got 0")
		}

		if count > 0 {
			return validateFirstPublicNetworkEssentials(rs.Primary.Attributes)
		}

		return nil
	}
}

func validateFirstPublicNetworkEssentials(attributes map[string]string) error {
	idKey := "public_networks.0.id"
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("first public network missing ID")
	}

	publicNameKey := "public_networks.0.public_name"
	if _, exists := attributes[publicNameKey]; !exists {
		return fmt.Errorf("first public network missing public_name")
	}

	datacenterIdKey := "public_networks.0.datacenter_id"
	if _, exists := attributes[datacenterIdKey]; !exists {
		return fmt.Errorf("first public network missing datacenter_id")
	}

	startDateKey := "public_networks.0.start_date"
	if _, exists := attributes[startDateKey]; !exists {
		return fmt.Errorf("first public network missing start_date")
	}

	sameVlanKey := "public_networks.0.same_vlan"
	if sameVlan, exists := attributes[sameVlanKey]; !exists {
		return fmt.Errorf("first public network missing same_vlan")
	} else if sameVlan != "true" && sameVlan != "false" {
		return fmt.Errorf("first public network same_vlan field must be 'true' or 'false', got '%s'", sameVlan)
	}

	typeKey := "public_networks.0.type"
	if _, exists := attributes[typeKey]; !exists {
		return fmt.Errorf("first public network missing type")
	}

	stateKey := "public_networks.0.state"
	if _, exists := attributes[stateKey]; !exists {
		return fmt.Errorf("first public network missing state")
	}

	serversKey := "public_networks.0.servers.#"
	if _, exists := attributes[serversKey]; !exists {
		return fmt.Errorf("first public network missing servers list")
	}

	ipsKey := "public_networks.0.ips.#"
	if _, exists := attributes[ipsKey]; !exists {
		return fmt.Errorf("first public network missing ips list")
	}

	availabilityZonesKey := "public_networks.0.availability_zones.#"
	if _, exists := attributes[availabilityZonesKey]; !exists {
		return fmt.Errorf("first public network missing availability_zones list")
	}

	lastLogsKey := "public_networks.0.last_logs.#"
	if _, exists := attributes[lastLogsKey]; !exists {
		return fmt.Errorf("first public network missing last_logs list")
	}

	return nil
}

func testAccPublicNetworksResourceConfig() string {
	return `
resource "arsys-baremetal_public_network" "test" {
  public_name   = "test-public-network"
  description   = "Test public network for acceptance testing"
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
}
`
}

func testAccPublicNetworksDataSourceConfig() string {
	return `
data "arsys-baremetal_public_networks" "all" {
}
`
}
