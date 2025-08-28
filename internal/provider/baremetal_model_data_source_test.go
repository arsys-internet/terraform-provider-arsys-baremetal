package provider

import (
	"fmt"
	"regexp"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccBaremetalModelDataSource(t *testing.T) {
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
				Config: testAccBaremetalModelDataSourceConfig("D603ACCF941E23C57B9653E619F1061E"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("state_id"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("hardware"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("hardware").AtMapKey("core"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("hardware").AtMapKey("ram"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("hardware").AtMapKey("hdds"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("hardware").AtMapKey("hdds").AtSliceIndex(0),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("hardware").AtMapKey("hdds").AtSliceIndex(0).AtMapKey("size"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("availability"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("availability").AtSliceIndex(0),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_model.test",
						tfjsonpath.New("availability").AtSliceIndex(0).AtMapKey("datacenter_id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
				},
			},
		},
	})
}
func testAccBaremetalModelDataSourceConfig(id string) string {
	return fmt.Sprintf(`
data "arsys-baremetal_baremetal_model" "test" {
  id = "%s"
}
`, id)
}
