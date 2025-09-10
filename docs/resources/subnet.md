---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Subnet Resource"
sidebar_current: "docs-resource-subnet"
description: |-
  Creates and manages Subnets in Arsys Baremetal.
---

# arsys-baremetal\_subnet

Creates and manages **Subnets** in Arsys Baremetal.

## Example Usage

### Create a subnet

```hcl
resource "arsys-baremetal_subnet" "example" {
  mask          = 28
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
}
```

## Argument Reference

The following arguments are supported:
* *`mask`* - (Required) The subnet mask(24,25,26,27,28)
* *`datacenter_id`* - (Required) ID of the data center where the subnet will be created. Cannot be updated after
creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the subnet.
* `name` - The automatically generated name of the subnet.
* `description` - The description of the subnet.
* `mask` - The subnet mask(24,25,26,27,28).
* `datacenter_id` - ID of the data center where the subnet is located.
* `public_name` - The public name of the subnet.
* `ip` - The IP address of the subnet.
* `gateway` - The gateway IP address.
* `broadcast` - The broadcast IP address.
* `type_id` - The subnet type ID.
* `type` - The subnet type.
* `dhcp` - DHCP configuration status.
* `state_id` - The state ID of the subnet.
* `state` - The current state of the subnet.
* `start_date` - The creation date of the subnet.
* `description` - The description of the subnet (if available).
* `subnet` - The subnet range (if available).
* `subnet_id` - The parent subnet ID (if available).
* `network_interface_id` - The network interface ID (if available).
* `server_type` - The server type (if available).
* `load_balancer_id` - The load balancer ID (if available).
* `load_balancer_public_name` - The load balancer public name (if available).
* `network_id` - The network ID (if available).
* `network_public_name` - The network public name (if available).
* `inversedns` - The inverse DNS configuration (if available).
* `last_logs` - List of recent log entries for the subnet.

## Notes

* **No Import Support**: This resource cannot be imported. Existing subnets must be managed outside Terraform or
  recreated.
* **No Update Support**: Subnets cannot be updated after creation. Any configuration changes require destroying and
  recreating the resource.
* **Immutable Configuration**: Both `mask` and `datacenter_id` cannot be changed after creation.

## Destroy

To destroy a subnet, use the following command:

```shell
terraform destroy -target=arsys-baremetal_subnet.example
```