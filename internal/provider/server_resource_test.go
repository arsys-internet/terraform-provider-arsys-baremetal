package provider

import (
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccServerResource(t *testing.T) {
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
			// Step 1: Create with full configuration
			{
				Config: testAccServerResourceFullConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test-baremetal-full"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Full configuration baremetal server"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("appliance_id"),
						knownvalue.StringExact("6EE23B88AB3CBA9944334C06E9075061"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("baremetal_model_id"),
						knownvalue.StringExact("650D003D3FC8A8FE554330E869B39FC0"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("power_on"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("install_backup_agent"),
						knownvalue.Bool(false),
					),

					// Verify that the server is created with a valid ID and type
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("server_type"),
						knownvalue.StringExact("baremetal"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("status").AtMapKey("state"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("vcore"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("ram"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("ips"),
						knownvalue.ListSizeExact(1),
					),
				},
			},
			// Step 2: Update name and description
			{
				Config: testAccServerResourceUpdatedConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test-baremetal-updated"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated description for baremetal server"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("baremetal_model_id"),
						knownvalue.StringExact("650D003D3FC8A8FE554330E869B39FC0"),
					),
				},
			},
			{
				ResourceName:      "arsys-baremetal_server.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password", "appliance_id", "datacenter_id",
				},
			},
		},
	})
}

func testAccServerResourceUpdatedConfig() string {
	return `
resource "arsys-baremetal_server" "test" {
  name           = "test-baremetal-updated"
  description    = "Updated description for baremetal server"
  appliance_id   = "6EE23B88AB3CBA9944334C06E9075061"
  datacenter_id  = "81DEF28500FBC2A973FC0C620DF5B721"
  power_on                = false
  install_backup_agent    = false
  
  hardware = {
    baremetal_model_id = "650D003D3FC8A8FE554330E869B39FC0"
  }
}
`
}

func testAccServerResourceFullConfig() string {
	return `
resource "arsys-baremetal_server" "test" {
  name           = "test-baremetal-full"
  description    = "Full configuration baremetal server"
  appliance_id   = "6EE23B88AB3CBA9944334C06E9075061"
  datacenter_id  = "81DEF28500FBC2A973FC0C620DF5B721"
  
  # Optional configuration fields
  power_on                = false
  install_backup_agent    = false
  
  hardware = {
    baremetal_model_id = "650D003D3FC8A8FE554330E869B39FC0"
  }
}
`
}
