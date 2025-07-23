package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"strconv"
	"testing"
	"time"
)

func TestAccServersDataSource(t *testing.T) {
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
				Config: testAccServersDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_servers.all",
						tfjsonpath.New("id"),
						knownvalue.StringExact("servers"),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_servers.all",
						tfjsonpath.New("servers"),
						knownvalue.NotNull(),
					),
				},
				Check: resource.ComposeTestCheckFunc(
					testAccCheckServersDataSourceConditional(),
				),
			},
		},
	})
}

func testAccCheckServersDataSourceConditional() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_servers.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr := rs.Primary.Attributes["servers.#"]
		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("servers count is not a number: %v", err)
		}

		if count == 0 {
			return nil
		}

		if serverType := rs.Primary.Attributes["servers.0.server_type"]; serverType != "baremetal" {
			return fmt.Errorf("expected server_type 'baremetal', got '%s'", serverType)
		}

		requiredFields := []string{"id", "name", "datacenter.id", "hardware.vcore", "status.state"}
		for _, field := range requiredFields {
			key := fmt.Sprintf("servers.0.%s", field)
			if value := rs.Primary.Attributes[key]; value == "" {
				return fmt.Errorf("missing required field: %s", field)
			}
		}

		return nil
	}
}

func testAccServersDataSourceConfig() string {
	return `
data "arsys-baremetal_servers" "all" {
}
`
}
