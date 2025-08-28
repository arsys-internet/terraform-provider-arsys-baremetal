---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal : Subnet Data source"
sidebar_current: "docs-datasource-subnet"
description: |-
  Get information about Subnets in Arsys Baremetal.
---

# arsys-baremetal\_subnets

The **Subnets data source** can be used to search and return all existing subnets.
Also, it can be used to search for and return an existing subnet.
You can provide a string for the id parameter which will be compared with the created subnet.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all subnets

```hcl
data "arsys-baremetal_public_subnets" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `subnets` - List of subnets

### Get By Id

```hcl
data "arsys-baremetal_public_subnet" "example" {
  id = subnet_id
}
```

## Argument Reference

* `id` - (Required)Id of an existing subnet that you want to search for.
  `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the subnet
* `ip` - The address of the subnet
* `type` - The type of subnet ("IPV4Subnet," "IPV6Subnet")
* `assigned_to` - Public network assigned to the subnet
    * `id` - Identifier of the resource
    * `name` - The name of the resource
    * `type` - The type of the resource ("PUBLIC_NETWORK")
* `subnet_id` - Subnet id
* `reverse_dns` - The reverse dns of the subnet
* `is_dhcp` - True if subnet is assigned using DHCP
* `state` - The state of the subnet
* `datacenter` - The data center where the subnet is located
    * `id` - Identifier of the Data Center
    * `location` - The regional location where the Data Center will be created
    * `country_code` - The country code of the Data Center
* `creation_date` - The creation date of the subnet