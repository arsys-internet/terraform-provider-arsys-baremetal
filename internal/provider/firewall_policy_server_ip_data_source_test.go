package provider

//import (
//	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
//	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
//	"github.com/hashicorp/terraform-plugin-testing/statecheck"
//	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
//	"github.com/hashicorp/terraform-plugin-testing/tfversion"
//	"regexp"
//	"terraform-provider-arsys-baremetal/internal/util"
//	"testing"
//	"time"
//)
// TODO: implementar el test de data source para server ip de firewall policy
//func TestAccFirewallPolicyServerIpDataSource(t *testing.T) {
//	resource.Test(t, resource.TestCase{
//		PreCheck: func() {
//			TestAccPreCheck(t)
//			time.Sleep(5 * time.Second)
//		},
//		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//			tfversion.SkipBelow(tfversion.Version1_0_0),
//		},
//		Steps: []resource.TestStep{
//			{
//				Config: testAccFirewallPolicyServerIPDataSourceConfig(),
//				// Sin ConfigStateChecks - no valida nada
//			},
//		},
//	})
//}
//
//func testAccFirewallPolicyServerIPDataSourceConfig() string {
//	return `
//data "arsys-baremetal_firewall_policy_server_ip" "test" {
//  firewall_policy_id = "983B8140A859CF5ACA8784D2B9ECC95F"
//  server_ip_id       = "08D42710D6071AE597C46EF3C6EB2272"
//}
//`
//}
