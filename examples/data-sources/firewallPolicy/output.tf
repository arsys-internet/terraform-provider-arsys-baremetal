//Output for all firewall policies
output "all_firewall_policies" {
    value = data.arsys-baremetal_firewall_policies.all
}
// Output for a specific firewall policy by its ID
output "firewall_policy" {
  value = data.arsys-baremetal_firewall_policy.example
}