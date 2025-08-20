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

func TestAccDatacentersDataSource(t *testing.T) {
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
				Config: testAccDatacentersDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_datacenters.all", "id"),
					testAccCheckDatacentersListNotEmpty(),
				),
			},
		},
	})
}

func testAccCheckDatacentersListNotEmpty() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_datacenters.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["datacenters.#"]
		if !ok {
			return fmt.Errorf("datacenters count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid datacenters count '%s': %v", countStr, err)
		}

		if count == 0 {
			return fmt.Errorf("expected at least 1 datacenter, got 0")
		}

		if count > 0 {
			return validateFirstDatacenterEssentials(rs.Primary.Attributes)
		}

		return nil
	}
}

func validateFirstDatacenterEssentials(attributes map[string]string) error {
	idKey := "datacenters.0.id"
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("first datacenter missing ID")
	}

	countryCodeKey := "datacenters.0.country_code"
	if _, exists := attributes[countryCodeKey]; !exists {
		return fmt.Errorf("first datacenter missing country_code")
	}

	locationKey := "datacenters.0.location"
	if _, exists := attributes[locationKey]; !exists {
		return fmt.Errorf("first datacenter missing location")
	}

	defaultKey := "datacenters.0.default"
	if defaultVal, exists := attributes[defaultKey]; !exists {
		return fmt.Errorf("first datacenter missing default")
	} else if defaultVal != "1" && defaultVal != "0" {
		return fmt.Errorf("first datacenter default field must be '1' or '0', got '%s'", defaultVal)
	}

	return nil
}

func testAccDatacentersDataSourceConfig() string {
	return `
data "arsys-baremetal_datacenters" "all" {
}
`
}
