---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Firewall Policy Data Sources"
sidebar_current: "docs-datasource-firewall_policy"
description: |-
  Get information about Firewall Policies in Arsys Baremetal
---

# arsys-baremetal\_firewall\_policies

The **Firewall policies data source** can be used to search and return all existing firewall policies.
Also, it can be used to search for and return an existing firewall policy.
You can provide a string for the id parameter which will be compared with the created firewall policy.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all Firewall Policies

```hcl
data "arsys-baremetal_firewall_policies" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `firewall_policies` - List of firewall policies

### Get By ID

```hcl
data "arsys-baremetal_firewall_policy" "example" {
  id = firewall_policy_id
}
```

## Argument Reference

* `id` - (Required) ID of an existing firewall policy that you want to search for.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the Firewall Policy
* `name` - Name of the Firewall Policy
* `description` - Description of the Firewall Policy
* `state` - State of the Firewall Policy (e.g., active, inactive)
* `creation_date` - Date when the Firewall Policy was created
* `default` - Indicates if the Firewall Policy is the default one in the panel
* `rules` - List of rules associated with the Firewall Policy, each rule containing:
    * `id` - Identifier of the rule
    * `protocol` - Protocol of the rule (e.g., TCP, UDP)
    * `port_from` - Starting port of the rule
    * `port_to` - Ending port of the rule (if applicable)
    * `source` - Source IP or CIDR of the rule
    * `description` - Description of the rule
    * `action` - Action of the rule (e.g., allow, deny)
* `server_ips` - List of server IPs associated with the Firewall Policy, each containing:
    * `id` - Identifier of the server IP
    * `ip` - IP address of the server
    * `server_name` - Name of the server associated with the IP
* `cloudpanel_id` - Cloud panel identifier of the Firewall Policy