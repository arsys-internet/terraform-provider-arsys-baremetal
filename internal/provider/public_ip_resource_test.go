package provider

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
	"regexp"
	"strings"
	service "terraform-provider-arsys-baremetal/internal/services/publicIp"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"
	"time"
)

func TestAccPublicIpResource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			TestAccPreCheck(t)
			time.Sleep(5 * time.Second)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		CheckDestroy: testAccCheckPublicIpDestroy,
		Steps: []resource.TestStep{
			// Test Create
			{
				Config: testAccPublicIpResourceConfig(
					"81DEF28500FBC2A973FC0C620DF5B721",
					"IPV4",
					"test-reverse-dns.example.com",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("IPV4"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("reverse_dns"),
						knownvalue.StringExact("test-reverse-dns.example.com"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("ip"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("is_dhcp"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("datacenter").AtMapKey("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("assigned_to"),
						knownvalue.Null(),
					),
				},
			},
			// Test Import
			{
				ResourceName:      "arsys-baremetal_public_ip.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"reverse_dns",
				},
			},
			// Test Update
			{
				Config: testAccPublicIpResourceConfig(
					"81DEF28500FBC2A973FC0C620DF5B721",
					"IPV4",
					"updated-reverse-dns.example.com",
				),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact("81DEF28500FBC2A973FC0C620DF5B721"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("type"),
						knownvalue.StringExact("IPV4"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("reverse_dns"),
						knownvalue.StringExact("updated-reverse-dns.example.com"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("ip"),
						knownvalue.StringRegexp(regexp.MustCompile(util.IPv4Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("state"),
						knownvalue.StringExact("ACTIVE"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("creation_date"),
						knownvalue.StringRegexp(regexp.MustCompile(util.DateTimePattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_public_ip.test",
						tfjsonpath.New("datacenter"),
						knownvalue.NotNull(),
					),
				},
			},
		},
	})
}

func testAccCheckPublicIpDestroy(s *terraform.State) error {
	testService := getTestPublicIpService()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "arsys-baremetal_public_ip" {
			continue
		}

		id := rs.Primary.ID
		if id == "" {
			continue
		}

		_, err := testService.GetPublicIp(id)

		if err == nil {
			return fmt.Errorf("public ip %s still exists", id)
		}

		if strings.Contains(err.Error(), "not found") {
			continue
		}

		return fmt.Errorf("error checking public ip %s: %s", id, err)
	}

	return nil
}

func testAccPublicIpResourceConfig(datacenterID, ipType, reverseDns string) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_public_ip" "test" {
  datacenter_id = %[1]q
  type          = %[2]q
  reverse_dns   = %[3]q
}
`, datacenterID, ipType, reverseDns)
}

func getTestPublicIpService() *service.ApiPublicIpService {
	apiClient := getTestClient()
	return service.NewApiPublicIpService(apiClient)
}
