---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Public Network IPs Resource"
sidebar_current: "docs-resource-public_network_ips"
description: |-
  Associate and disassociate IPs with Public Networks in Arsys Baremetal.
---

# arsys-baremetal\_public\_network\_ips

Associate and disassociate **IPs** with **Public Networks** in Arsys Baremetal.

## Example Usage

### Associate IPs with public network.

```hcl
resource "arsys-baremetal_public_network_ips" "example" {
  public_network_id = "F17227DA6D4883B596A3861AA994EA56"
  action            = true
  ips = ["08CB1405D8732066D531486E7B7AAC31"]
}
```

## Argument Reference

The following arguments are supported:

* `public_network_id` - (Required) The ID of the public network to which the IPs will be associated.
* `action` - (Required) Whether to associate or disassociate the IPs.
* `ips` - (Required) A list of IPs IDs to associate with the public network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - A list of all IPs associating with the public network. Each item contains the following attributes:
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

### Disassociate IPs with public network.

```hcl
resource "arsys-baremetal_public_network_ips" "example" {
  public_network_id = "F17227DA6D4883B596A3861AA994EA56"
  action            = false
  ips = ["08CB1405D8732066D531486E7B7AAC31"]
}
```

## Argument Reference

The following arguments are supported:

* `public_network_id` - (Required) The ID of the public network to which the IPs will be associated.
* `action` - (Required) Whether to associate or disassociate the IPs.
* `ips` - (Required) A list of IPs IDs to associate with the public network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `items` - A list of all IPs disassociating with the public network. Each item contains the following attributes:
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
