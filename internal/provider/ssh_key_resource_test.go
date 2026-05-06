package provider

import (
	"errors"
	"fmt"
	"regexp"
	service "terraform-provider-arsys-baremetal/internal/services/sshkey"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccSshKeyResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
			time.Sleep(5 * time.Second)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		CheckDestroy: testAccCheckSshKeyDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccSshKeyResourceConfig(
					"acc-test-ssh-key-terraform",
					"SSH key for acceptance testing",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("acc-test-ssh-key-terraform"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("SSH key for acceptance testing"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("public_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("md5"),
						knownvalue.StringRegexp(regexp.MustCompile(util.Md5Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("servers"),
						knownvalue.ListSizeExact(0),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("private_key"),
						knownvalue.NotNull(),
					),
				},
			},
			{
				ResourceName:      "arsys-baremetal_ssh_key.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"description",
					"private_key",
				},
			},
			{
				Config: testAccSshKeyResourceConfig(
					"acc-test-ssh-key-terraform-updated",
					"Updated SSH key description for acceptance testing",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("acc-test-ssh-key-terraform-updated"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated SSH key description for acceptance testing"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("public_key"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_ssh_key.test",
						tfjsonpath.New("md5"),
						knownvalue.StringRegexp(regexp.MustCompile(util.Md5Pattern)),
					),
				},
			},
		},
	})
}

func testAccCheckSshKeyDestroy(s *terraform.State) error {
	testService := getTestSshKeyService()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "arsys-baremetal_ssh_key" {
			continue
		}

		id := rs.Primary.ID
		if id == "" {
			continue
		}

		_, err := testService.GetSshKey(id)

		if err == nil {
			return fmt.Errorf("ssh key %s still exists", id)
		}

		if errors.Is(err, util.ErrNotFound) {
			continue
		}

		return fmt.Errorf("error checking ssh key %s: %s", id, err)
	}

	return nil
}

func testAccSshKeyResourceConfig(name string, description string) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_ssh_key" "test" {
  name        = %[1]q
  description = %[2]q
}
`, name, description)
}

func getTestSshKeyService() *service.ApiSshKeyService {
	apiClient := getTestClient()
	return service.NewApiSshKeyService(apiClient)
}
