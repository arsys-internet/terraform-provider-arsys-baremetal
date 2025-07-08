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
				Config: testAccPublicIpDataSourceConfig("D9C3C724CC521387E962E822F043017D"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_public_ip.test",
						tfjsonpath.New("id"),
						knownvalue.StringExact("D9C3C724CC521387E962E822F043017D"),
					),
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
						knownvalue.NotNull(),
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
						tfjsonpath.New("assigned_to"),
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

func testAccPublicIpDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "arsys-baremetal_public_ip" "test" {
  id = "%s"
}
`, id)
}
