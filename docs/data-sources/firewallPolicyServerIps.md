---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Firewall Policy Server IPs Data Sources"
sidebar_current: "docs-datasource-firewall_policy_server_ips"
description: |-
  Get information about Firewall Policy Server Ips of a Firewall Policy in Arsys Baremetal
---

# arsys-baremetal\_firewall\_policy\_server\_ips

The **Firewall policy server ips data source** can be used to search and return all existing server ips from a specific
firewall policy.
Also, it can be used to search for and return a specific server ip from a firewall policy.
You can provide a string for the id parameter which will be compared with the existing firewall policy server ips.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all Server Ips from a Firewall Policy

```hcl
data "arsys-baremetal_firewall_policy_server_ips" "example" {
  id = "983B8140A859CF5ACA8784D2B9ECC95F"
}
```

## Argument Reference

* `id` - (Required) ID of an existing firewall policy that you want to retrieve server ips from.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the Firewall Policy
* `server_ips` - List of server IPs associated with the Firewall Policy, each server IP containing:
    * `id` - Identifier of the server IP
    * `ip` - IP address of the server
    * `server_name` - Name of the server

### Get Specific Server Ip by Id

## Example Usage

```hcl
data "arsys-baremetal_firewall_policy_server_ip" "example" {
  firewall_policy_id = "983B8140A859CF5ACA8784D2B9ECC95F"
  server_ip_id       = "08D42710D6071AE597C46EF3C6EB2272"
}
```

## Argument Reference

* `firewall_policy_id` - (Required) ID of an existing firewall policy containing the server ip.
* `server_ip_id` - (Required) ID of the specific server ip that you want to retrieve.

Both `firewall_policy_id` and `server_ip_id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the server IP
* `firewall_policy_id` - Id of the parent firewall policy
* `ip` - IP address of the server
* `server_name` - Name of the server