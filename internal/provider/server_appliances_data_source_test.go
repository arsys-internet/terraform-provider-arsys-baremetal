package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"strconv"
	"testing"
	"time"
)

func TestAccServerAppliancesDataSource(t *testing.T) {
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
				Config: testAccServerAppliancesDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_server_appliances.all", "id"),
					testAccCheckServerAppliancesListNotEmpty(),
				),
			},
		},
	})
}

func testAccCheckServerAppliancesListNotEmpty() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_server_appliances.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["server_appliances.#"]
		if !ok {
			return fmt.Errorf("server_appliances count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid server_appliances count '%s': %v", countStr, err)
		}

		if count == 0 {
			return fmt.Errorf("expected at least 1 server appliance, got 0")
		}

		if count > 0 {
			return validateFirstServerApplianceEssentials(rs.Primary.Attributes)
		}

		return nil
	}
}

func validateFirstServerApplianceEssentials(attributes map[string]string) error {
	idKey := "server_appliances.0.id"
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("first server appliance missing ID")
	}

	nameKey := "server_appliances.0.name"
	if _, exists := attributes[nameKey]; !exists {
		return fmt.Errorf("first server appliance missing name")
	}

	osFamilyKey := "server_appliances.0.os_family"
	if _, exists := attributes[osFamilyKey]; !exists {
		return fmt.Errorf("first server appliance missing os_family")
	}

	osKey := "server_appliances.0.os"
	if _, exists := attributes[osKey]; !exists {
		return fmt.Errorf("first server appliance missing os")
	}

	osVersionKey := "server_appliances.0.os_version"
	if _, exists := attributes[osVersionKey]; !exists {
		return fmt.Errorf("first server appliance missing os_version")
	}

	osArchKey := "server_appliances.0.os_architecture"
	if _, exists := attributes[osArchKey]; !exists {
		return fmt.Errorf("first server appliance missing os_architecture")
	}

	typeKey := "server_appliances.0.type"
	if _, exists := attributes[typeKey]; !exists {
		return fmt.Errorf("first server appliance missing type")
	}

	minHddSizeKey := "server_appliances.0.min_hdd_size"
	if _, exists := attributes[minHddSizeKey]; !exists {
		return fmt.Errorf("first server appliance missing min_hdd_size")
	}

	versionKey := "server_appliances.0.version"
	if _, exists := attributes[versionKey]; !exists {
		return fmt.Errorf("first server appliance missing version")
	}

	dcCountKey := "server_appliances.0.available_datacenters.#"
	if dcCount, exists := attributes[dcCountKey]; !exists || dcCount == "0" {
		return fmt.Errorf("first server appliance must have at least one available datacenter")
	}

	compatCountKey := "server_appliances.0.server_type_compatibility.#"
	if compatCount, exists := attributes[compatCountKey]; !exists || compatCount == "0" {
		return fmt.Errorf("first server appliance must have at least one server type compatibility")
	}

	categoriesCountKey := "server_appliances.0.categories.#"
	if categoriesCount, exists := attributes[categoriesCountKey]; !exists || categoriesCount == "0" {
		return fmt.Errorf("first server appliance must have at least one category")
	}

	return nil
}

func testAccServerAppliancesDataSourceConfig() string {
	return `
data "arsys-baremetal_server_appliances" "all" {
  # Este data source lista todos los server appliances
}
`
}
