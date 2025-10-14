package provider

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccServerApplianceDataSource(t *testing.T) {
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
				Config: testAccServerApplianceDataSourceConfig("65754F5F403868DB67EDD90CC204C1BE"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("65754F5F403868DB67EDD90CC204C1BE"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("available_datacenters"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("os_image_type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("os_family"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("os"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("os_version"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("os_architecture"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("min_hdd_size"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("licenses"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("type"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("server_type_compatibility"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server_appliance.test",
						tfjsonpath.New("categories"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccServerApplianceDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "arsys-baremetal_server_appliance" "test" {
  id = "%s"
}
`, id)
}
