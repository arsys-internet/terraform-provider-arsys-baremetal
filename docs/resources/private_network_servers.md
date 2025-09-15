---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Private Network Server Resources"
sidebar_current: "docs-resource-private_network_servers"
description: |-
  Manages servers for an existing Private Network in Arsys Baremetal
---

# arsys-baremetal\_private\_network\_servers\_assign

Assign servers to an existing **Private Network** in Arsys baremetal.

## Example Usage

### Assigns servers to an existing Private Network

```hcl
resource "arsys-baremetal_private_network_servers_assign" "example" {
  id = "6EAC18BDB1084D6C2F5C984DE1E3EDBB"
  servers = ["825CD55B22A61C75A9B9ED48EC80D2EE", "A982AE3D56CEB4CEB2FB9B62C6B74691"]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The ID of the existing Private Network to add servers to
* `servers` - (Required) A list of servers to be added to the Private Network. Each element is an id of a server
  that you want to assign to the Private Network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

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
* `cloudpanel_id` - CloudPanel ID of the private network

## Notes

* This resource cannot be updated. Any changes require replacement.
* At least one server is required.

---


Removes a specific server from an existing **Private Network** in Arsys baremetal.

## Example Usage

### Remove a server from an existing Private Network

```hcl
resource "arsys-baremetal_private_network_server_remove" "example" {
  id        = "4EFAD5836CE43ACA502FD5B99BEE44EF"
  server_id = "8FA2D5836CE43ACA502FD5B99BEE77AB"
}
```

## Argument Reference

The following arguments are supported:

* `id`- (Required) The ID of the existing private network to remove the server from.
* `server_id` - (Required) The ID of the specific server to remove from the private network.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

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

## Import

This resource cannot be imported as it represents an operation to remove a server from a private network.

## Notes

* This resource cannot be updated. Any changes require replacement.
* The delete operation only removes the resource from the Terraform state.

# Destroy Examples

## Destroy private network server assign resource

```shell
terraform destroy -target=arsys-baremetal_private_network_servers_assign.example
```

## Destroy a private network server remove resource

```shell
terraform destroy -target=arsys-baremetal_private_network_server_remove.example
```