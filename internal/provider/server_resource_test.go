package provider

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"terraform-provider-arsys-baremetal/internal/services/datacenter"
	"terraform-provider-arsys-baremetal/internal/services/server"
	"terraform-provider-arsys-baremetal/internal/services/serverappliance"
	"terraform-provider-arsys-baremetal/internal/util"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

type testServerIds struct {
	modelId      string
	datacenterId string
	applianceId  string
}

func pickAvailableServerTestIds(t *testing.T) testServerIds {
	t.Helper()
	c := getTestClient()
	serverSvc := server.NewApiServerService(c)
	applianceSvc := serverappliance.NewApiServerApplianceService(c)
	datacenterSvc := datacenter.NewApiDatacenterService(c)

	allDatacenters, err := datacenterSvc.GetDatacenters()
	if err != nil {
		t.Skipf("skipping: could not retrieve datacenters: %s", err)
	}
	spainDCs := make(map[string]struct{})
	for _, dc := range allDatacenters {
		if strings.EqualFold(dc.CountryCode, "ES") {
			spainDCs[dc.Id] = struct{}{}
		}
	}
	if len(spainDCs) == 0 {
		t.Skip("skipping: no Spanish datacenters found")
	}

	baremetalModels, err := serverSvc.GetBaremetalModels()
	if err != nil {
		t.Skipf("skipping: could not retrieve baremetal models: %s", err)
	}

	type candidate struct {
		modelId      string
		datacenterId string
		ram          int
		cores        int
	}
	var candidates []candidate
	for _, m := range baremetalModels {
		for _, a := range m.Availability {
			if _, inSpain := spainDCs[a.DatacenterId]; inSpain && a.Available {
				candidates = append(candidates, candidate{
					modelId:      m.Id,
					datacenterId: a.DatacenterId,
					ram:          m.Hardware.Ram,
					cores:        m.Hardware.Core,
				})
			}
		}
	}
	if len(candidates) == 0 {
		t.Skip("skipping: no baremetal model available in Spanish datacenters")
	}
	sort.Slice(candidates, func(i, j int) bool {
		if candidates[i].ram != candidates[j].ram {
			return candidates[i].ram < candidates[j].ram
		}
		return candidates[i].cores < candidates[j].cores
	})
	chosen := candidates[0]

	appliances, err := applianceSvc.GetServerAppliances()
	if err != nil {
		t.Skipf("skipping: could not retrieve server appliances: %s", err)
	}
	var applianceId string
	for _, a := range appliances {
		if !strings.Contains(strings.ToLower(a.Name), "ubuntu") {
			continue
		}
		for _, dc := range a.AvailableDatacenters {
			if dc == chosen.datacenterId {
				applianceId = a.Id
				break
			}
		}
		if applianceId != "" {
			break
		}
	}
	if applianceId == "" {
		t.Skipf("skipping: no Ubuntu appliance available in datacenter %s", chosen.datacenterId)
	}

	return testServerIds{
		modelId:      chosen.modelId,
		datacenterId: chosen.datacenterId,
		applianceId:  applianceId,
	}
}

func TestAccServerResource(t *testing.T) {
	TestAccPreCheck(t)
	ids := pickAvailableServerTestIds(t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			time.Sleep(5 * time.Second)
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_0_0),
		},
		Steps: []resource.TestStep{
			{
				Config: testAccServerResourceFullConfig(ids),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test-baremetal-full"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Full configuration baremetal server"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("appliance_id"),
						knownvalue.StringExact(ids.applianceId),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("datacenter_id"),
						knownvalue.StringExact(ids.datacenterId),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("baremetal_model_id"),
						knownvalue.StringExact(ids.modelId),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("power_on"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("install_backup_agent"),
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("server_type"),
						knownvalue.StringExact("baremetal"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("status").AtMapKey("state"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("vcore"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("ram"),
						knownvalue.NotNull(),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("ips"),
						knownvalue.ListSizeExact(1),
					),
				},
			},
			{
				Config: testAccServerResourceUpdatedConfig(ids),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("name"),
						knownvalue.StringExact("test-baremetal-updated"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("description"),
						knownvalue.StringExact("Updated description for baremetal server"),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("id"),
						knownvalue.StringRegexp(regexp.MustCompile(util.HexID32Pattern)),
					),
					statecheck.ExpectKnownValue(
						"arsys-baremetal_server.test",
						tfjsonpath.New("hardware").AtMapKey("baremetal_model_id"),
						knownvalue.StringExact(ids.modelId),
					),
				},
			},
			{
				ResourceName:      "arsys-baremetal_server.test",
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"password",
					"appliance_id",
					"datacenter_id",
					"power_on",
					"install_backup_agent",
				},
			},
		},
	})
}

func testAccServerResourceUpdatedConfig(ids testServerIds) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_server" "test" {
 name           = "test-baremetal-updated"
 description    = "Updated description for baremetal server"
 appliance_id   = %q
 datacenter_id  = %q
 power_on                = false
 install_backup_agent    = false

 hardware = {
   baremetal_model_id = %q
 }
}
`, ids.applianceId, ids.datacenterId, ids.modelId)
}

func testAccServerResourceFullConfig(ids testServerIds) string {
	return fmt.Sprintf(`
resource "arsys-baremetal_server" "test" {
 name           = "test-baremetal-full"
 description    = "Full configuration baremetal server"
 appliance_id   = %q
 datacenter_id  = %q
 power_on                = false
 install_backup_agent    = false

 hardware = {
   baremetal_model_id = %q
 }
}
`, ids.applianceId, ids.datacenterId, ids.modelId)
}
