package provider

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccSshKeysDataSource(t *testing.T) {
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
				Config: testAccSshKeysResourceConfig(),
			},
			{
				Config: testAccSshKeysDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.arsys-baremetal_ssh_keys.all", "id"),
					testAccCheckSshKeysListNotEmpty(),
				),
			},
		},
	})
}

func testAccCheckSshKeysListNotEmpty() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources["data.arsys-baremetal_ssh_keys.all"]
		if !ok {
			return fmt.Errorf("data source not found")
		}

		countStr, ok := rs.Primary.Attributes["ssh_keys.#"]
		if !ok {
			return fmt.Errorf("ssh_keys count not found")
		}

		count, err := strconv.Atoi(countStr)
		if err != nil {
			return fmt.Errorf("invalid ssh_keys count '%s': %v", countStr, err)
		}

		if count == 0 {
			return fmt.Errorf("expected at least 1 SSH key, got 0")
		}

		if count > 0 {
			return validateFirstSshKeyEssentials(rs.Primary.Attributes)
		}

		return nil
	}
}

func validateFirstSshKeyEssentials(attributes map[string]string) error {
	idKey := "ssh_keys.0.id"
	if id, exists := attributes[idKey]; !exists || id == "" {
		return fmt.Errorf("first SSH key missing ID")
	}

	nameKey := "ssh_keys.0.name"
	if _, exists := attributes[nameKey]; !exists {
		return fmt.Errorf("first SSH key missing name")
	}

	stateKey := "ssh_keys.0.state"
	if _, exists := attributes[stateKey]; !exists {
		return fmt.Errorf("first SSH key missing state")
	}

	serversKey := "ssh_keys.0.servers.#"
	if _, exists := attributes[serversKey]; !exists {
		return fmt.Errorf("first SSH key missing servers list")
	}

	md5Key := "ssh_keys.0.md5"
	if md5, exists := attributes[md5Key]; !exists || md5 == "" {
		return fmt.Errorf("first SSH key missing md5")
	}

	publicKeyKey := "ssh_keys.0.public_key"
	if publicKey, exists := attributes[publicKeyKey]; !exists || publicKey == "" {
		return fmt.Errorf("first SSH key missing public_key")
	}

	creationDateKey := "ssh_keys.0.creation_date"
	if _, exists := attributes[creationDateKey]; !exists {
		return fmt.Errorf("first SSH key missing creation_date")
	}

	return nil
}

func testAccSshKeysResourceConfig() string {
	return `
resource "arsys-baremetal_ssh_key" "test" {
  name        = "ssh-key-list-test-terraform"
  description = "SSH key for list testing"
}
`
}

func testAccSshKeysDataSourceConfig() string {
	return testAccSshKeysResourceConfig() + `
data "arsys-baremetal_ssh_keys" "all" {
}
`
}
