---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal : Public ip Data source"
sidebar_current: "docs-datasource-public_ip"
description: |-
  Get information on a Arsys Baremetal public ip
---

# arsys-baremetal\_ip

The **Public IPs data source** can be used to search and return all existing public IPs.
Also, can be used to search for and return an existing IP.
You can provide a string for the id parameter which will be compared with created IPs.
If a single match is found, it will be returned. If it is not found an error will be returned.

## Example Usage

### Get all IPs

```hcl
data "arsys-baremetal_public_ips" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `ips` - List of IPs
  Algo

### Get By ID

```hcl
data "arsys-baremetal_public_ip" "example" {
  id = ip_id
}
```

## Argument Reference

* `id` - (Required)ID of an existing IP that you want to search for.
  `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the private network
* `ip` - The name of the private network
* `assigned_to` - Resource where the IP is assigned
    * `id` - Identifier of the resource
    * `name` - The name of the resource
    * `type` - The type of the resource ("SERVER", "LOAD_BALANCER")
* `subnet_id` - Subnet where the IP is located
* `reverse_dns` - The reverse dns of the IP
* `is_dhcp` - True if IP is assigned using DHCP
* `state` - The state of the IP
* `datacenter` - The data center where the IP is located
    * `id` - Identifier of the Data Center
    * `location` - The regional location where the Data Center will be created
    * `country_code` - The country code of the Data Center
* `creation_date` - The creation date of the IP