---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Firewall Policy Resources"
sidebar_current: "docs-resource-firewall_policy"
description: |-
  Creates and manages Firewall Policies in Arsys Baremetal
---

# arsys-baremetal\_firewall\_policy

Creates and manages **Firewall Policies** in Arsys baremetal.

## Example Usage

# Create a Firewall Policy.

```hcl
resource "arsys-baremetal_firewall_policy" "example" {
  name        = "Firewall Policy example"
  description = "Firewall Policy description"
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
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the Firewall Policy.
* `description` - (Optional) The description of the Firewall Policy.
* `rules` - (Optional) A list of rules to be applied to the Firewall Policy. Each rule should be a map containing:
    * `protocol` - (Required) The protocol of the rule (e.g., TCP, UDP).
    * `port_from` - (Required) The starting port of the rule.
    * `port_to` - (Required) The ending port of the rule (if applicable).
    * `source` - (Optional) The source IP or CIDR of the rule.

# Update a Firewall Policy.

```hcl
resource "arsys-baremetal_firewall_policy" "example_update" {
  name        = "New Firewall Policy example"
  description = "New Firewall Policy description"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Optional) The name of the Firewall Policy.
* `description` - (Optional) The description of the Firewall Policy.

# Import a Firewall Policy

Resource Firewall Policy can be imported using creating the resource the `resource ID`, e.g.
**Note:** The resource must be declared before importing. No arguments are required!.

Example:

```hcl
resource "arsys-baremetal_firewall_policy" "example_import" {}
```

```shell
terraform import arsys-baremetal_firewall_policy.example_import {firewall_policy id}
```

# Destroy a Firewall Policy

To destroy a Firewall Policy, use the following command:

```shell
terraform destroy -target=arsys-baremetal_firewall_policy.example
```