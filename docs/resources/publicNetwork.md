---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Public Network Resource"
sidebar_current: "docs-resource-public_network"
description: |-
  Creates and manages Public Networks in Arsys Baremetal.
---

# arsys-baremetal\_public\_network

Creates and manages **Public Networks** in Arsys Baremetal.

## Example Usage

### Create or update a public network.

```hcl
resource "arsys-baremetal_public_network" "example" {
  name          = "Public Network example"
  datacenter_id = "99DEF28511HBC2A973HC0C620DH5B732"
  description   = "public network description"
}
```

When updating a public network, datacenter_id is only needed for the creation.

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the public network.
* `datacenter_id` - (Required) The ID of the data center where the public network will be created.
* `description` - (Optional) The description of the public network.

### Import

Resource Public Network can be imported using creating the resource the `resource id`, e.g.
**Note:** The resource must be declared before importing. No arguments are required!

Example:

```hcl
resource "arsys-baremetal_public_network" "example_import" {}
```

```shell
terraform import arsys-baremetal_public_network.example_import {publicNetwork uuid}
```

### Destroy

To destroy a public network, use the following command:

```shell
terraform destroy -target=arsys-baremetal_public_network.example
```