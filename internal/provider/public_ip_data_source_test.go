package provider

import (
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccPublicIpDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccPublicIpResourceConfig(
					"81DEF28500FBC2A973FC0C620DF5B721",
					"IPV4",
					"test-datasource.example.com",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccPublicIpDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("ip"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("IPV4"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("is_dhcp"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
				},
			},
		},
	})
}

func testAccPublicIpDataSourceConfig() string {
	return testAccPublicIpResourceConfig("81DEF28500FBC2A973FC0C620DF5B721", "IPV4", "test-datasource.example.com") + `
data "arsys-baremetal_public_ip" "test" {
  id = arsys-baremetal_public_ip.test.id
}
`
}
