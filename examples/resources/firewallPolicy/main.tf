//Resource that creates a firewall policy in a baremetal environment using Terraform
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
    //Additional rules can be added here
  ]
}

// Resource that updates an existing firewall policy
resource "arsys-baremetal_firewall_policy" "test_policy_update" {
  name        = "terraform_test_updated"
  description = "update"
}

# Resource to import an existing firewall policy by its ID
resource "arsys-baremetal_firewall_policy" "test_policy_import" {
}
// Execute terraform import command to import the existing firewall policy