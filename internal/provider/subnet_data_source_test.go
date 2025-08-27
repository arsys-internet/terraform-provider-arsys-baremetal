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

func TestAccSubnetDataSource(t *testing.T) {
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
				Config: testAccSubnetResourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_subnet.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccSubnetDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_subnet.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_subnet.test",
						tfjsonpath.New("ip"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_subnet.test",
						tfjsonpath.New("type"),
						knownvalue.StringRegexp(regexp.MustCompile("^IPV[46]Subnet$")),
					),
				},
			},
		},
	})
}

func testAccSubnetDataSourceConfig() string {
	return testAccSubnetResourceConfig() + `
data "arsys-baremetal_subnet" "test" {
  id = arsys-baremetal_subnet.test.id
}
`
}
