---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Private Network Servers Data Sources"
sidebar_current: "docs-datasource-private_network_servers"
description: |-
  Get information about Servers assigned to a Private Network in Arsys Baremetal
---

# arsys-baremetal\_private\_network\_servers

The **Private network servers data source** can be used to search and return all servers assigned to a specific
private network.
Also, it can be used to search for and return a specific server assigned to a private network.
You can provide a string for the id parameter which will be compared with the existing private network server.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all servers assigned to a Private Network

```hcl
data "arsys-baremetal_private_network_servers" "all" {
  id = "983B8140A859CF5ACA8784D2B9ECC95F"
}
```

## Argument Reference

* `id` - (Required) ID of an existing private network that you want to retrieve servers from.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the Private Network
* `servers` - List of servers assigned to the Private Network, each server containing:
    * `id` - Identifier of the server
    * `name` - Name of the server
    * `lock` - Status of the server lock (integer)

### Get Specific Server by ID

## Example Usage

```hcl
data "arsys-baremetal_private_network_server" "example" {
  private_network_id = "983B8140A859CF5ACA8784D2B9ECC95F"
  id                 = "08D42710D6071AE597C46EF3C6EB2272"
}
```

## Argument Reference

* `private_network_id` - (Required) ID of an existing private network containing the server ip.
* `id` - (Required) ID of the specific server ip that you want to retrieve.

Both `private_network_id` and `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:
* `private_network_id` - Identifier of private network.
* `id` - Identifier of the server
* `name` - Name of the server
* `lock` - Status of the server lock (integer)