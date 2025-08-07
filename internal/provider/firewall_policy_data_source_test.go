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

func TestAccFirewallPolicyDataSource(t *testing.T) {
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
				Config: testAccFirewallPolicyDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("firewall-policy-with-rules"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_firewall_policy.test",
						tfjsonpath.New("rules"),
						knownvalue.SetSizeExact(2),
					),
				},
			},
		},
	})
}

func testAccFirewallPolicyDataSourceConfig() string {
	return testFirewallPolicyWithRulesConfig + `

data "arsys-baremetal_firewall_policy" "test" {
  id = arsys-baremetal_firewall_policy.test.id
}
`
}
