---
subcategory: "Security"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: SSH Key Data Source"
sidebar_current: "docs-data-source-ssh_key"
description: |-
  Provides information about an SSH Keys in Arsys Baremetal.
---

# arsys-baremetal\_ssh\_key

Provides information about an **SSH Key** in Arsys Baremetal.

## Example Usage

### Get all SSH keys

```hcl
data "arsys-baremetal_ssh_keys" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `ssh_keys` - List of SSH keys

### Get By ID

```hcl
data "arsys-baremetal_ssh_key" "example" {
  id = ssh_key_id
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) ID of an existing SSH key that you want to search for.

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
* `creation_date` - SSH key creation date in ISO 8601 format (e.g., 2023-05-29T09:43:31+00:00).
* `private_key` - The SSH private key (if available).

