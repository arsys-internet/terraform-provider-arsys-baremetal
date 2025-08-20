---
subcategory: "Compute Engine"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Server"
sidebar_current: "docs-resource-server"
description: |-
  Creates and manages Servers in Arsys Baremetal.
---

# arsys-baremetal_server

Creates and manages **Servers** in Arsys Baremetal.

## Example Usage

The following examples show how to create or update a server.

### Basic Configuration (Minimal Required Fields)

```hcl
resource "arsys-baremetal_server" "basic_server" {
  name = "baremetal-server"
  appliance_id = "6EE23B88AB3CBA9944334C06E9075061"  # Rocky Linux 9
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"  # Spain datacenter
  hardware = {
    baremetal_model_id = "650D003D3FC8A8FE554330E869B39FC0"
  }
}
```

### Complete Configuration (All Available Options)

```hcl
resource "arsys-baremetal_server" "complete_server" {
  name        = "baremetal-server"
  description = "Baremetal server for testing"
  appliance_id = "6EE23B88AB3CBA9944334C06E9075061"  # Rocky Linux 9
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"  # Spain datacenter
  hardware = {
    baremetal_model_id = "650D003D3FC8A8FE554330E869B39FC0"
  }
  password = "MySecurePassword123!"
  power_on = true
  install_backup_agent = true

  # Network and security configuration  
  firewall_policy_id   = "A1B2C3D4E5F6789012345678901234AB"
  ip_id                = "B2C3D4E5F6789012345678901234ABCD"
  load_balancer_id     = "C3D4E5F6789012345678901234ABCDEF"
  monitoring_policy_id = "D4E5F6789012345678901234ABCDEF12"
  availability_zone_id = "E5F6789012345678901234ABCDEF1234"
}
```

### Update Server


```hcl
resource "arsys-baremetal_server" "example" {
  name        = "baremetal-server-updated"
  description = "Updated server description"
}
```

## Argument Reference

For the creation of a baremetal server, the following arguments are supported:

### Required Arguments

* `name` - (Required) The name of the server
* `appliance_id` - (Required) The ID of the OS appliance/image to install
* `datacenter_id` - (Required) The ID of the datacenter where the server will be created
* `hardware` - (Required) Hardware configuration block
    * `baremetal_model_id` - (Required) The ID of the baremetal model to use

### Optional Arguments

* `description` - The description of the server
* `password` - Password for the server. If not provided, a random password will be generated
* `power_on` - (Boolean) Whether to power on the server after creation. Defaults to `true`
* `firewall_policy_id` - The ID of the firewall policy to associate with the server
* `ip_id` - The ID of an existing IP to assign to the server
* `load_balancer_id` - The ID of a load balancer to associate with the server
* `monitoring_policy_id` - The ID of a monitoring policy to associate with the server
* `install_backup_agent` - (Boolean) Whether to install the backup agent. Defaults to `false`
* `availability_zone_id` - The ID of the availability zone

## Import

Server resources can be imported using the server ID:

**Note:** The resource must be declared before importing. No arguments are required for import.

Example:

```hcl
resource "arsys-baremetal_server" "example_server_import" {}
```

```shell
terraform import arsys-baremetal_server.example_server_import {server_id}
```

## Destroy

Server resources can be destroyed using the resource name and target option:

```shell
terraform destroy -target=arsys-baremetal_server.complete_server
```