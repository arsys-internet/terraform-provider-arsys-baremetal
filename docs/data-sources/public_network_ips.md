---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Public Network IPs Data Sources"
sidebar_current: "docs-datasource-public_network_ips"
description: |-
  Get information about IPs in a Public Networks in Arsys Baremetal.
---

# arsys-baremetal\_public\_network\_ips

The **Public Network IPs data source** can be used to search and return all existing IPS in a public network.
Also, it can be used to search for and return an existing IP in a public network.
You can provide a string for the id parameter which will be compared with created public networks.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all IPs in a Public Network

```hcl
data "arsys-baremetal_public_network_ips" "all" {
  public_network_id = public_network_id
}
```

## Argument Reference

* `public_network_id` - (Required) ID of an existing public network that you want to search for.

`public_network_id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `public_network_ips` - List of IPs in a public network

### Get IP By ID

```hcl
data "arsys-baremetal_public_network_ip" "example" {
  public_network_id = public_network_id
}
```

## Argument Reference

* `public_network_id` - (Required) ID of an existing public network that you want to search for.

`public_network_id` must be provided. If none, the datasource will return an error.

* `id` - (Required) ID of an existing IP in the public network that you want to search for.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of IP in the public network
* `ip_address` - The IP address of the public network
* `description` - The description of the IP in the public network
* `network_interface_id` - The network interface ID where the IP is located
* `lb_id` - The load balancer ID where the IP is located
* `inverse_dns` - Inverse DNS of the IP in the public network
* `start_date` - Date when the IP in the public network was created
* `site_id` - ID of the site where the IP in the public network is located
* `is_main` - Indicates if the IP is the main one in the public network
* `mask` - The mask of the IP in the public network
* `firewall_id` - ID of the firewall associated with the IP in the public network
* `gateway` - The gateway of IP in the public network
* `broadcast` - The broadcast address of the IP in the public network
* `network_id` - ID of the network associated with the IP in the public network
* `type` - The type of the IP in the public network
* `state` - The state of the IP in the public network