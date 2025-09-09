---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Public Network IPs Resource"
sidebar_current: "docs-resource-public_network_ips"
description: |-
  Associate and disassociate IPs with Public Networks in Arsys Baremetal.
---

# arsys-baremetal\_public\_network\_ips

Associate and disassociate **Servers** with **Public Networks** in Arsys Baremetal.

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
