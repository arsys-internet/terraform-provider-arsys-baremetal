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

### Get By ID

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

* `id` - Identifier of the private network
* `name` - The name of the private network
* `description` - The description of the private network
* `network_address` - The network address of the private network
* `subnet_mask` - The subnet mask of the private network
* `state` - The state of the private network
* `datacenter` - The data center where the private network is located
    * `id` - Identifier of the Data Center
    * `country_code` - The country code of the Data Center
    * `location` - The regional location where the Data Center will be created
* `creation_date` - The creation date of the private network
* `servers` - The servers that are connected to the private network
    * `id` - Identifier of the server
    * `name` - The name of the server
* `cloudpanel_id` - The CloudPanel ID of the private network
