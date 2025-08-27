---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal : Public ip Data source"
sidebar_current: "docs-datasource-public_ip"
description: |-
  Creates and manages Public IPs in Arsys Baremetal .
---

# arsys-baremetal\_public\_ip

Creates and manages **Public IPs** in Arsys Baremetal.

## Example Usage

### Create or update an IP.

```hcl
resource "arsys-baremetal_public_ip" "example" {
  reverse_dns   = "dns123.domain.com"
  datacenter_id = "99DEF28511HBC2A973HC0C620DH5B732"
  type          = "IPV4"
}
```

## Argument Reference

The following arguments are supported:

* `reverse_dns` - (Optional) The reverse dns name of the IP.
* `datacenter_id` - (_Only for create_, Optional) Identifier of the datacenter where the IP will be created (only for
  unassigned IPs)
* `type` - (_Only for create_, Optional) The type of the IP ("IPV4", "IPV6").

### Import

Resource IP can be imported using creating the resource the `resource id`, e.g.
**Note:** The resource must be declared before importing. No arguments are required!.

Example:

```hcl
resource "arsys-baremetal_public_ip" "example_import" {}
```

```shell
terraform import arsys-baremetal_public_ip.example_import {ip uuid}
```

### Destroy

To destroy an IP, use the following command:

```shell
terraform destroy -target=arsys-baremetal_public_ip.example