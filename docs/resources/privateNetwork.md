---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Private Network Resource"
sidebar_current: "docs-resource-privateNetwork"
description: |-
  Creates and manages Private Networks in Arsys Baremetal .
---

# arsys-baremetal\_private_network

Creates and manages **Private Networks** in Arsys Baremetal.

## Example Usage

### Create or update a private network.

```hcl
resource "arsys-baremetal_private_network" "example" {
  name            = "Private Network example"
  datacenter_id   = "99DEF28511HBC2A973HC0C620DH5B732"
  description     = "private network description"
  network_address = "192.168.1.0"
  subnet_mask     = "255.255.255.0"
}
```

When updating a private network, datacenter_id is only needed for the creation.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the private network.
* `datacenter_id` - (Required) The ID of the data center where the private network will be created.
* `description` - (Optional) The description of the private network.
* `network_address` - (Optional) The network address of the private network.
* `subnet_mask` - (Optional) The subnet mask of the private network.

### Import

Resource Private Network can be imported using creating the resource the `resource id`, e.g.
**Note:** The resource must be declared before importing. No arguments are required!

Example:

```hcl
resource "arsys-baremetal_private_network" "example_import" {}
```

```shell
terraform import arsys-baremetal_private_network.example_import {privateNetwork uuid}
```

### Destroy

To destroy a private network, use the following command:

```shell
terraform destroy -target=arsys-baremetal_private_network.example
```