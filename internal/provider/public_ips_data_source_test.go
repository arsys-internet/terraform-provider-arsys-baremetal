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

func TestAccPublicIpsDataSource(t *testing.T) {
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
				Config: testAccPublicIpsResourceConfig(),
			},
			{
				Config: testAccPublicIpsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_public_ips.all", "id"),
					testAccCheckPublicIpsListNotEmpty(),
				),
			},
		},
	})
}

func testAccCheckPublicIpsListNotEmpty() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_public_ips.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["public_ips.#"]
		if !ok {
			return fmt.Errorf("public_ips count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid public_ips count '%s': %v", countStr, err)
		}

		if count == 0 {
			return fmt.Errorf("expected at least 1 public IP, got 0")
		}

		if count > 0 {
			return validateFirstPublicIpEssentials(rs.Primary.Attributes)
		}

		return nil
	}
}

func validateFirstPublicIpEssentials(attributes map[string]string) error {
	idKey := "public_ips.0.id"
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("first public IP missing ID")
	}

	ipKey := "public_ips.0.ip"
	if _, exists := attributes[ipKey]; !exists {
		return fmt.Errorf("first public IP missing ip")
	}

	typeKey := "public_ips.0.type"
	if _, exists := attributes[typeKey]; !exists {
		return fmt.Errorf("first public IP missing type")
	}

	stateKey := "public_ips.0.state"
	if _, exists := attributes[stateKey]; !exists {
		return fmt.Errorf("first public IP missing state")
	}

	isDhcpKey := "public_ips.0.is_dhcp"
	if isDhcp, exists := attributes[isDhcpKey]; !exists {
		return fmt.Errorf("first public IP missing is_dhcp")
	} else if isDhcp != "true" && isDhcp != "false" {
		return fmt.Errorf("first public IP is_dhcp field must be 'true' or 'false', got '%s'", isDhcp)
	}

	if assignedToId, exists := attributes["public_ips.0.assigned_to.id"]; exists && assignedToId != "" {
		assignedToIdKey := "public_ips.0.assigned_to.id"
		if _, exists := attributes[assignedToIdKey]; !exists {
			return fmt.Errorf("first public IP missing assigned_to.id")
		}
	}

	datacenterIdKey := "public_ips.0.datacenter.id"
	if _, exists := attributes[datacenterIdKey]; !exists {
		return fmt.Errorf("first public IP missing datacenter.id")
	}

	creationDateKey := "public_ips.0.creation_date"
	if _, exists := attributes[creationDateKey]; !exists {
		return fmt.Errorf("first public IP missing creation_date")
	}

	return nil
}

func testAccPublicIpsResourceConfig() string {
	return `
resource "arsys-baremetal_public_ip" "test" {
   reverse_dns = "dns1324.dominio.com"
   datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
   type = "IPV4"
}
`
}

func testAccPublicIpsDataSourceConfig() string {
	return `
data "arsys-baremetal_public_ips" "all" {
}
`
}
