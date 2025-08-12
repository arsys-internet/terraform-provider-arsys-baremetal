package provider

import (
	"fmt"
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

func TestAccFirewallPolicyResource(t *testing.T) {
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
				Config: testFirewallPolicyWithRulesConfig,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("firewall-policy-with-rules"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("rules"),
						knownvalue.SetSizeExact(2),
					),
				},
			},
			{
				ResourceName:      "arsys-baremetal_firewall_policy.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccFirewallPolicyResourceConfig(
					"firewall-policy-updated",
					"Updated description",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("firewall-policy-updated"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated description"),
					),
				},
			},
		},
	})
}

func testAccFirewallPolicyResourceConfig(name, description string) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_firewall_policy" "test" {
  name        = "%s"
  description = "%s"
}
`, name, description)
}

var testFirewallPolicyWithRulesConfig = `
resource "arsys-baremetal_firewall_policy" "test" {
  name        = "firewall-policy-with-rules"
  description = "Firewall policy with rules for testing"
  rules = [
    {
      protocol  = "TCP"
      port_from = 22
      port_to   = 22
      source    = "192.168.1.0/24"
    },
    {
      protocol  = "TCP"
      port_from = 80
      port_to   = 80
      source    = "0.0.0.0"
    }
  ]
}
`
