package provider

import (
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
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("first_password"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hostname"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("server_type"),
						knownvalue.StringExact("baremetal"),
					),

					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("managed"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("ssh_password"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("rsa_key"),
						knownvalue.Bool(false),
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

					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("image"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware"),
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
						tfjsonpath.New("hardware").AtMapKey("hdds"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("status"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("status").AtMapKey("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("connection_speed"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("redundancy"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("ips"),
						knownvalue.ListSizeExact(1),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("private_networks"),
						knownvalue.ListSizeExact(0),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("password"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("firewall_policy_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("ip_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("load_balancer_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("monitoring_policy_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("availability_zone_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("dvd"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("alerts"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("monitoring_policy"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("cloudpanel_id"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("snapshot"),
						knownvalue.Null(),
					),

					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("recovery_mode"),
						knownvalue.StringRegexp(regexp.MustCompile("^(false|true)$")),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("recovery_image_os"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("recovery_user"),
						knownvalue.Null(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("recovery_password"),
						knownvalue.Null(),
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
  
  # Keep the same optional configuration as Step 1
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
