---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Public Network Servers Resource"
sidebar_current: "docs-resource-public_network_servers"
description: |-
  Associate and disassociate Servers with Public Networks in Arsys Baremetal.
---

# arsys-baremetal\_public\_network\_servers

Associate and disassociate **Servers** with **Public Networks** in Arsys Baremetal.

## Example Usage

### Associate and disassociate servers with public network.

```hcl
resource "arsys-baremetal_public_network_servers" "example" {
  public_network_id = "F17227DA6D4883B596A3861AA994EA56"
  servers = ["08CB1405D8732066D531486E7B7AAC31"]
}
```

## Argument Reference

The following arguments are supported:

* `public_network_id` - (Required) ID of the public network to which the servers will be associated.
* `servers` - (Required) A list of server IDs to associate with the public network. To disassociate a server, remove it
  from the list.
