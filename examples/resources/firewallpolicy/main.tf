# Example that creates a firewall policy in a baremetal environment using Terraform
resource "arsys-baremetal_firewall_policy" "test_policy" {
  name        = "terraform_test_policy"
  description = "Firewall policy created with Terraform for test"
  rules = [
    {
      protocol  = "TCP"
      port_from = 22
      port_to   = 22
      source    = "192.168.1.0/24"
    },
    # Additional rules can be added here
  ]
}

# Example that updates an existing firewall policy
resource "arsys-baremetal_firewall_policy" "test_policy" {
  name        = "terraform_test_updated"
  description = "update"
}

# Example to import an existing firewall policy by ID
resource "arsys-baremetal_firewall_policy" "test_policy_import" {
}
# Execute terraform import command to import the existing firewall policy


# Example that adds a rule to an existing firewall policy
resource "arsys-baremetal_firewall_policy_rule_add" "allow_udp" {
  id = var.firewall_policy_id
  rules = [
    {
      protocol  = "UDP"
      port_from = 85
      port_to   = 85
      source    = "8.8.8.8"
      action    = "allow"
      description = "Allow DNS queries from Public DNS"
    }
    # Additional rules can be added here
  ]
}

# Example that removes a rule from an existing firewall policy
resource "arsys-baremetal_firewall_policy_rule_remove" "remove" {
  id = var.firewall_policy_id
  rule_id = var.rule_id
}

# Example that assigns server IPs to a firewall policy
resource "arsys-baremetal_firewall_policy_server_ips" "example" {
  id = var.firewall_policy_id
  server_ips = ["08D42710D6071AE597C46EF3C6EB2272"] # Additional server IP Ids can be added inside the list
}