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

func TestAccBaremetalModelsDataSource(t *testing.T) {
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
				Config: testAccBaremetalModelsDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("state_id"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("hardware"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("hardware").AtMapKey("core"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("hardware").AtMapKey("ram"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("hardware").AtMapKey("hdds"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("hardware").AtMapKey("hdds").AtSliceIndex(0),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("hardware").AtMapKey("hdds").AtSliceIndex(0).AtMapKey("size"),
						knownvalue.NotNull(),
					),

					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("availability"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("availability").AtSliceIndex(0),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("availability").AtSliceIndex(0).AtMapKey("datacenter_id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_baremetal_models.test",
						tfjsonpath.New("baremetal_models").AtSliceIndex(0).AtMapKey("availability").AtSliceIndex(0).AtMapKey("available"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccBaremetalModelsDataSourceConfig() string {
	return `
data "arsys-baremetal_baremetal_models" "test" {
}
`
}
