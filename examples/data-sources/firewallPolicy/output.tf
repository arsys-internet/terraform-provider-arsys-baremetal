//Output for all firewall policies
output "all_firewall_policies" {
  value = data.arsys-baremetal_firewall_policies.all
}
// Output for a specific firewall policy by its ID
output "firewall_policy" {
  value = data.arsys-baremetal_firewall_policy.example
}

// Output for all firewall policy rules
output "firewall_policy_rules" {
  value = data.arsys-baremetal_firewall_policy_rules.example
}

// Output for a specific firewall policy rule by Id
output "firewall_policy_rule" {
  value = data.arsys-baremetal_firewall_policy_rule.example
}