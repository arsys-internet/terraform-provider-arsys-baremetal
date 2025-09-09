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

func TestAccSshKeyDataSource(t *testing.T) {
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
				Config: testAccSshKeyResourceConfig(
					"test-ssh-key",
					"Test SSH key for acceptance testing",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				Config: testAccSshKeyDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("name"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("description"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("servers"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("md5"),
						knownvalue.StringRegexp(regexp.MustCompile(util.Md5Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("public_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_ssh_key.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
				},
			},
		},
	})
}

func testAccSshKeyDataSourceConfig() string {
	return testAccSshKeyResourceConfig("test-ssh-key", "Test SSH key for acceptance testing") + `
data "arsys-baremetal_ssh_key" "test" {
  id = arsys-baremetal_ssh_key.test.id
}
`
}
