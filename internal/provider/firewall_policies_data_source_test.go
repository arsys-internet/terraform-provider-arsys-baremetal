package provider

import (
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccFirewallPoliciesDataSource(t *testing.T) {
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
				Config: testAccFirewallPoliciesDataSourceConfig(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_firewall_policies.all",
						tfjsonpath.New("firewall_policies"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"data.arsys-baremetal_firewall_policies.all",
						tfjsonpath.New("id"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccFirewallPoliciesDataSourceConfig() string {
	return `
data "arsys-baremetal_firewall_policies" "all" {}
`
}
