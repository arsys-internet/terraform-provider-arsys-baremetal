# Firewall Policy

This example demonstrates how to create and manage firewall policies using the Arsys Baremetal provider.

## What This Example Does

- Creates a firewall policy using `arsys-baremetal_firewall_policy`
- Updates the firewall policy `arsys-baremetal_firewall_policy`
- You can also delete the firewall policy using `terraform destroy -target=arsys-baremetal_firewall_policy.name_of_policy`
- Import an existing firewall policy using `terraform import arsys-baremetal_firewall_policy.name_of_policy <policy_id>`
- Add a firewall rule to the firewall policy using `arsys-baremetal_firewall_rule_add`
- Remove a firewall rule from the firewall policy using `arsys-baremetal_firewall_rule_remove`