# Output for all firewall policies
output "all_firewall_policies" {
  value = data.arsys-baremetal_firewall_policies.all
}
# Output for a specific firewall policy by its ID
output "firewall_policy" {
  value = data.arsys-baremetal_firewall_policy.example
}

# Output for all firewall policy rules
output "firewall_policy_rules" {
  value = data.arsys-baremetal_firewall_policy_rules.example
}

# Output for a specific firewall policy rule by Id
output "firewall_policy_rule" {
  value = data.arsys-baremetal_firewall_policy_rule.example
}

# Output for all server IPs assigned to a firewall policy
output "firewall_policy_server_ips" {
  value = data.arsys-baremetal_firewall_policy_server_ips.all
}