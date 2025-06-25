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

func TestAccPrivateNetworkDataSource(t *testing.T) {
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
				Config: testAccPrivateNetworkResourceConfig(
					"pn-datasource-test",
					"81DEF28500FBC2A973FC0C620DF5B721",
					"10.0.0.0",
					"255.255.255.0",
					"Private network for datasource testing"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_private_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("pn-datasource-test"),
					),
				},
			},
			{
				Config: testAccPrivateNetworkDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("pn-datasource-test"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("name"),
						knownvalue.StringRegexp(regexp.MustCompile(util.NamePattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Private network for datasource testing"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("network_address"),
						knownvalue.StringExact("10.0.0.0"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("network_address"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("subnet_mask"),
						knownvalue.StringExact("255.255.255.0"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("subnet_mask"),
						knownvalue.StringRegexp(regexp.MustCompile(util.SubnetMaskPattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("cloudpanel_id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("datacenter").AtMapKey("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_private_network.test",
						tfjsonpath.New("servers"),
						knownvalue.ListSizeExact(0),
					),
				},
			},
		},
	})
}

func testAccPrivateNetworkDataSourceConfig() string {
	return testAccPrivateNetworkResourceConfig("pn-datasource-test", "81DEF28500FBC2A973FC0C620DF5B721", "10.0.0.0", "255.255.255.0", "Private network for datasource testing") + `
data "arsys-baremetal_private_network" "test" {
  id = arsys-baremetal_private_network.test.id
}
`
}
