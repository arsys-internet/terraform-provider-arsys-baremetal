# Output for a specific firewall policy by ID
output "firewall_policy_resource" {
  value = arsys-baremetal_firewall_policy.test_policy
}

# Output for a specific firewall policy after a rule is added
output "firewall_policy_rule_add" {
  value = arsys-baremetal_firewall_policy_rule_add.allow_udp
}

# Output for a specific firewall policy after a rule is removed
output "firewall_policy_server_rule_remove" {
  value = arsys-baremetal_firewall_policy_rule_remove.remove
}

# Output for a specific firewall policy after a server IP is assigned
output "firewall_policy_server_ips" {
  value = arsys-baremetal_firewall_policy_server_ips.example
}