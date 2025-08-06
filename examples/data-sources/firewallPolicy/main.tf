//Example to discover all firewall policies that have been created
data "arsys-baremetal_firewall_policies" "all" {}

//Example to retrieve a specific firewall policy by ID
data "arsys-baremetal_firewall_policy" "example" {
  id = "firewall_policy_id" # Replace with the actual firewall policy ID
}