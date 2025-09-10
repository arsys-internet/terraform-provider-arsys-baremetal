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

func TestAccSubnetResource(t *testing.T) {
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
				Config: testAccSubnetResourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("mask"),
						knownvalue.Int64Exact(28),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("ip"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("public_name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("gateway"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("broadcast"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
				},
			},
			{
				Config:      testAccSubnetResourceUpdatedConfig(),
				ExpectError: regexp.MustCompile("Update not supported"),
			},
		},
	})
}

func testAccSubnetResourceConfig() string {
	return `
resource "arsys-baremetal_subnet" "test" {
  mask          = 28
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
}`
}

func testAccSubnetResourceUpdatedConfig() string {
	return `
resource "arsys-baremetal_subnet" "test" {
  mask          = 27
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
}`
}
