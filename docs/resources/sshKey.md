---
subcategory: "Security"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: SSH Key Resource"
sidebar_current: "docs-resource-ssh_key"
description: |-
  Creates and manages SSH Keys in Arsys Baremetal.
---

# arsys-baremetal\_ssh\_key

Creates and manages **SSH Keys** in Arsys Baremetal.

## Example Usage

### Create or update a SSH key.

```hcl 
resource "arsys-baremetal_ssh_key" "example" {
  name        = "my-ssh-key"
  description = "SSH key for accessing servers"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7..."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the SSH key.
* `description` - (Optional) The description of the SSH key.
* `public_key` - (_Only for create_, Optional) The public key content. If not provided, a new key pair will be
  generated.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - Identifier of the SSH key.
* `name` - The name of the SSH key.
* `description` - The description of the SSH key.
* `state` - Current state of the SSH key.
* `servers` - List of servers associated with the SSH key. Each item contains:
    * `id` - Server identifier.
    * `name` - Server name.
* `md5` - MD5 hash of the SSH key (32 hexadecimal characters).
* `public_key` - The SSH public key.
* `creation_date` - SSH key creation date in ISO 8601 format.
* `private_key` - SSH key private key (only available when SSH key is created).

### Import

Resource SSH Key can be imported using the `resource id`, e.g.
**Note:** The resource must be declared before importing. No arguments are required!

Example:

```hcl
resource "arsys-baremetal_ssh_key" "example_import" {}
```

```shell
terraform import arsys-baremetal_ssh_key.example_import {sshKey uuid}
```

### Destroy

To destroy an SSH key, use the following command:

```shell 
terraform destroy -target=arsys-baremetal_ssh_key.example
``` 

