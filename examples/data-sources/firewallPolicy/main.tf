# Example to discover all firewall policies that have been created
data "arsys-baremetal_firewall_policies" "all" {}

# Example to retrieve a specific firewall policy by ID
data "arsys-baremetal_firewall_policy" "example" {
  id = var.firewall_policy_id # Replace with the actual firewall policy ID
}

# Example to discover all firewall policy rules that have been created
data "arsys-baremetal_firewall_policy_rules" "example" {
  id = var.firewall_policy_id
}

# Example to retrieve a specific firewall policy rule by ID
data "arsys-baremetal_firewall_policy_rule" "example" {
  firewall_policy_id = var.firewall_policy_id
  id                 = var.rule_id
}

# Example to discover all server IPs that have been associated with a firewall policy
data "arsys-baremetal_firewall_policy_server_ips" "all" {
  id = var.firewall_policy_id
}

# Example to retrieve a specific server IP assigned to a firewall policy
data "arsys-baremetal_firewall_policy_server_ip" "example" {
  firewall_policy_id = var.firewall_policy_id
  server_ip_id       = var.server_ip_id
}
