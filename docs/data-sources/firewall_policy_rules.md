---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Firewall Policy Rules Data Sources"
sidebar_current: "docs-datasource-firewall_policy_rules"
description: |-
  Get information about Firewall Policy Rules in Arsys Baremetal
---

# arsys-baremetal\_firewall_policy_rules

The **Firewall policy rules data source** can be used to search and return all existing rules from a specific firewall
policy.
Also, it can be used to search for and return a specific rule from a firewall policy.
You can provide a string for the id parameter which will be compared with the existing firewall policy rules.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all Rules from a Firewall Policy

```hcl
data "arsys-baremetal_firewall_policy_rules" "all_rules" {
  id = "983B8140A859CF5ACA8784D2B9ECC95F"
}
```

## Argument Reference

* `id` - (Required) ID of an existing firewall policy that you want to retrieve rules from.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the Firewall Policy
* `rules` - List of rules associated with the Firewall Policy, each rule containing:
    * `id` - Identifier of the rule
    * `protocol` - Protocol of the rule (TCP, UDP, ICMP, TCP/UDP, IPSEC, GRE)
    * `port_from` - Starting port of the rule
    * `port_to` - Ending port of the rule
    * `source` - Source IP or CIDR of the rule
    * `description` - Description of the rule
    * `action` - Action of the rule (allow, deny)

### Get Specific Rule by ID

## Example Usage

```hcl
data "arsys-baremetal_firewall_policy_rule" "ssh_rule" {
  firewall_policy_id = "983B8140A859CF5ACA8784D2B9ECC95F"
  id                 = "4DA2DFD1E88347F24495F037490DB701"
}
```

## Argument Reference

* `firewall_policy_id` - (Required) ID of an existing firewall policy containing the rule.
* `id` - (Required) ID of the specific rule that you want to retrieve.

Both `firewall_policy_id` and `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the rule
* `firewall_policy_id` - ID of the parent firewall policy
* `protocol` - Protocol of the rule (TCP, UDP, ICMP, TCP/UDP, IPSEC, GRE)
* `port_from` - Starting port number
* `port_to` - Ending port number
* `source` - Source IP address or CIDR range
* `description` - Description of the rule
* `action` - Action performed by the rule (allow, deny)