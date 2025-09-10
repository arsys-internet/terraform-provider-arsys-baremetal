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

func TestAccSubnetsDataSource(t *testing.T) {
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
				Config: testAccSubnetsDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_subnets.all", "id"),
					testAccCheckSubnetsDataSourceValid(),
				),
			},
		},
	})
}

func testAccCheckSubnetsDataSourceValid() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_subnets.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["subnets.#"]
		if !ok {
			return fmt.Errorf("subnets count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid subnets count '%s': %v", countStr, err)
		}

		if count == 0 {
			return nil
		}

		return validateFirstSubnetEssentials(rs.Primary.Attributes)
	}
}

func validateFirstSubnetEssentials(attributes map[string]string) error {
	idKey := "subnets.0.id"
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("subnet missing ID")
	}

	ipKey := "subnets.0.ip"
	if _, exists := attributes[ipKey]; !exists {
		return fmt.Errorf("subnet missing ip")
	}

	typeKey := "subnets.0.type"
	subnetType, exists := attributes[typeKey]
	if !exists {
		return fmt.Errorf("subnet missing type")
	}

	if subnetType != "IPV4Subnet" && subnetType != "IPV6Subnet" {
		return fmt.Errorf("subnet type must be 'IPV4Subnet' or 'IPV6Subnet', got '%s'", subnetType)
	}

	datacenterIdKey := "subnets.0.datacenter.id"
	if _, exists := attributes[datacenterIdKey]; !exists {
		return fmt.Errorf("subnet missing datacenter.id")
	}

	return nil
}

func testAccSubnetsDataSourceConfig() string {
	return `data "arsys-baremetal_subnets" "all" {}`
}
