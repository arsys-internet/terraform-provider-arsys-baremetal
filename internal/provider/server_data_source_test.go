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
			// Step 1: Create the server resource
			{
				Config: testAccServerResourceConfig(
					"server-datasource-test",
					"Baremetal server created with Terraform for testing",
					"81DEF28500FBC2A973FC0C620DF5B721",
					"6EE23B88AB3CBA9944334C06E9075061",
					"650D003D3FC8A8FE554330E869B39FC0"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("server-datasource-test"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("status").AtMapKey("state"),
						knownvalue.NotNull(),
					),
				},
			},
			// Step 2: Use data source to read the created server
			{
				Config: testAccServerDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					// === CORE IDENTIFIERS ===
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
						tfjsonpath.New("description"),
						knownvalue.StringExact("Baremetal server created with Terraform for testing"),
					),

					// === SERVER PROPERTIES ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("server_type"),
						knownvalue.StringExact("baremetal"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("managed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hostname"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),

					// === AUTHENTICATION ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("ssh_password"),
						knownvalue.Bool(true), // Default para baremetal
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("rsa_key"),
						knownvalue.Bool(false), // Default para baremetal
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("first_password"),
						knownvalue.NotNull(), // Baremetal servers get a password
					),

					// === DATACENTER ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("datacenter").AtMapKey("id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("datacenter").AtMapKey("country_code"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("datacenter").AtMapKey("location"),
						knownvalue.NotNull(),
					),

					// === IMAGE ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("image"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("image").AtMapKey("id"),
						knownvalue.StringExact("6EE23B88AB3CBA9944334C06E9075061"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("image").AtMapKey("name"),
						knownvalue.NotNull(),
					),

					// === HARDWARE ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("baremetal_model_id"),
						knownvalue.StringExact("650D003D3FC8A8FE554330E869B39FC0"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("vcore"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("cores_per_processor"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("ram"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("hdds"),
						knownvalue.ListSizeExact(1), // Baremetal always has HDDs
					),

					// === STATUS ===
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
						tfjsonpath.New("status").AtMapKey("percent"),
						knownvalue.NotNull(),
					),

					// === NETWORKING ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("ips"),
						knownvalue.ListSizeExact(1), // Baremetal gets one IP
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("private_networks"),
						knownvalue.ListSizeExact(0), // No private networks by default
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("connection_speed"),
						knownvalue.NotNull(), // Baremetal has connection speed
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("redundancy"),
						knownvalue.NotNull(), // Baremetal has redundancy info
					),

					// === OPTIONAL FIELDS THAT SHOULD BE NULL ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("dvd"),
						knownvalue.Null(), // No DVD mounted by default
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("alerts"),
						knownvalue.Null(), // No alerts configured
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("monitoring_policy"),
						knownvalue.Null(), // No monitoring policy by default
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("snapshot"),
						knownvalue.Null(), // No snapshots by default
					),

					// === OPTIONAL FIELDS THAT MIGHT BE NULL ===
					// Note: hypervisor might be null for baremetal, so we don't check it
					// Note: cloudpanel_id might be null, so we don't check it

					// === VERIFY FIELD FORMATS ===
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringRegexp(regexp.MustCompile(util.NamePattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_server.test",
						tfjsonpath.New("datacenter").AtMapKey("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
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

func testAccServerDataSourceConfig() string {
	return testAccServerResourceConfig(
		"server-datasource-test",
		"Baremetal server created with Terraform for testing",
		"81DEF28500FBC2A973FC0C620DF5B721",
		"6EE23B88AB3CBA9944334C06E9075061",
		"650D003D3FC8A8FE554330E869B39FC0") + `

data "arsys-baremetal_server" "test" {
  id = arsys-baremetal_server.test.id
}
`
}
