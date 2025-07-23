package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"
	"time"
)

func TestAccServerDataSource(t *testing.T) {
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
				Config: testServerConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("server-datasource-test"),
					),
				},
			},
			{
				Config: testAccServerDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					// Verify that the data source retrieves the server details correctly
					//with the expected values only for essential attributes
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("server-datasource-test"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("server_type"),
						knownvalue.StringExact("baremetal"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("datacenter").AtMapKey("id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("image").AtMapKey("id"),
						knownvalue.StringExact("6EE23B88AB3CBA9944334C06E9075061"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("baremetal_model_id"),
						knownvalue.StringExact("650D003D3FC8A8FE554330E869B39FC0"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("vcore"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("status"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("status").AtMapKey("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
				},
			},
		},
	})
}

func testAccServerResourceConfig(name, description, datacenterId, applianceId, baremetalModelId string) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_server" "test" {
  name          = "%s"
  description   = "%s"
  datacenter_id = "%s"
  appliance_id  = "%s"
  hardware = {
    baremetal_model_id = "%s"
  }
}
`, name, description, datacenterId, applianceId, baremetalModelId)
}

var testServerConfig = testAccServerResourceConfig(
	"server-datasource-test",
	"Baremetal server created with Terraform for testing",
	"81DEF28500FBC2A973FC0C620DF5B721",
	"6EE23B88AB3CBA9944334C06E9075061",
	"650D003D3FC8A8FE554330E869B39FC0")

func testAccServerDataSourceConfig() string {
	return testServerConfig + `

data "arsys-baremetal_server" "test" {
  id = arsys-baremetal_server.test.id
}
`
}
