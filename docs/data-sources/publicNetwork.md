---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Public Network Data Sources"
sidebar_current: "docs-datasource-public_network"
description: |-
  Get information about Public Networks in Arsys Baremetal.
---

# arsys-baremetal\_public_network

The **Public Networks data source** can be used to search and return all existing public networks.
Also, it can be used to search for and return an existing public network.
You can provide a string for the id parameter which will be compared with created public networks.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all Public Networks

```hcl
data "arsys-baremetal_public_networks" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `public_networks` - List of public networks

### Get By Id

```hcl
data "arsys-baremetal_public_network" "example" {
  id = public_network_id
}
```

## Argument Reference

* `id` - (Required) ID of an existing public network that you want to search for.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the public network
* `public_name` - The name of the public network
* `description` - The description of the public network
* `datacenter_id` - The data center ID where the public network is located
* `start_date` - Date when the public network was created
* `same_vlan` - Indicates if the public network shares the same VLAN with other networks
* `type` - The type of the public network
* `state` - The state of the public network
* `servers` - The servers that are connected to the public network
    * `id` - Identifier of the server
    * `name` - The name of the server
    * `mac` - The MAC address of the server
    * `tagged` - Indicates if the server is tagged
* `ips` - List of IPs Id in the public network
* `availability_zones` - The availability zones of the public network
    * `id` - Identifier of the availability zone
    * `vlan_id` - Identifier of the VLAN
* `last_logs` - The last logs of the public network
    * `id` - Identifier of the log
    * `uuid` - The UUID of the log
    * `date` - The creation date of the operation
    * `action` - The action of the operation
    * `time` - The time it took for the operation
    * `result` - The result of the operation
    * `type` - The type of the operation
