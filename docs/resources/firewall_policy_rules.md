---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Firewall Policy Rule Resources"
sidebar_current: "docs-resource-firewall_policy_rules"
description: |-
  Manages rules for existing Firewall Policies in Arsys Baremetal
---

# arsys-baremetal\_firewall\_policy\_rule\_add

Adds rules to an existing **Firewall Policy** in Arsys baremetal.

## Example Usage

### Add rules to an existing Firewall Policy

```hcl
resource "arsys-baremetal_firewall_policy_rule_add" "example" {
  id = "4EFAD5836CE43ACA502FD5B99BEE44EF"
  rules = [
    {
      protocol  = "TCP"
      port_from = 80
      port_to   = 80
      source    = "0.0.0.0/0"
      action    = "allow"
      description = "Allow HTTP traffic"
    },
    {
      protocol  = "TCP"
      port_from = 443
      port_to   = 443
      source    = "0.0.0.0/0"
      action    = "allow"
      description = "Allow HTTPS traffic"
    }
  ]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) ID of the existing Firewall Policy to add rules to.
* `rules` - (Required) A list of rules to be added to the Firewall Policy. Each rule should contain:
    * `protocol` - (Required) The protocol of the rule (e.g., TCP, UDP).
    * `port_from` - (Required) The starting port of the rule.
    * `port_to` - (Required) The ending port of the rule.
    * `source` - (Optional) The source IP or CIDR of the rule.
    * `action` - (Optional) The action to take when the rule is matched (e.g., allow, deny)
    * `description` - (Optional) A description of the rule.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Firewall Policy.
* `name` - The name of the Firewall Policy.
* `description` - The description of the Firewall Policy.
* `state` - The current state of the Firewall Policy.
* `creation_date` - The date when the Firewall Policy was created.
* `default` - Define default panel firewalls.
* `cloudpanel_id` - Identifier of the cloud panel.
* `rules_detail` - Complete list of rules in the firewall policy after assignment.
* `server_ips` - Servers assigned to firewall policy.

## Notes

* This resource cannot be updated. Any changes require replacement.
* At least one rule is required.

---


Removes a specific rule from an existing **Firewall Policy** in Arsys baremetal.

## Example Usage

### Remove a rule from an existing Firewall Policy

```hcl
resource "arsys-baremetal_firewall_policy_rule_remove" "example" {
  id      = "4EFAD5836CE43ACA502FD5B99BEE44EF"
  rule_id = "8FA2D5836CE43ACA502FD5B99BEE77AB"
}
```

## Argument Reference

The following arguments are supported:

* `id`- (Required) ID of the existing Firewall Policy to remove the rule from.
* `rule_id` - (Required) ID of the specific rule to remove from the Firewall Policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - ID of the Firewall Policy.
* `name` - The name of the Firewall Policy.
* `description` - The description of the Firewall Policy.
* `state` - The current state of the Firewall Policy.
* `creation_date` - The date when the Firewall Policy was created.
* `default` - Define default panel firewalls.
* `rules` - Complete list of rules remaining in the firewall policy after removal.
* `server_ips` - Servers assigned to firewall policy.
* `cloudpanel_id` - Identifier of the cloud panel.

## Import

This resource cannot be imported as it represents an operation to remove a rule from an existing policy.

## Notes

* This resource cannot be updated. Any changes require replacement.
* The delete operation only removes the resource from the Terraform state.

# Destroy Examples

## Destroy firewall policy rule add resource

```shell
terraform destroy -target=arsys-baremetal_firewall_policy_rule_add.example
```

## Destroy firewall policy rule remove resource

```shell
terraform destroy -target=arsys-baremetal_firewall_policy_rule_remove.example
```