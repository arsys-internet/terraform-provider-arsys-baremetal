---
subcategory: "Network"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Firewall Policy Server IPs Resources"
sidebar_current: "docs-resource-firewall_policy_server_ips"
description: |-
  Manages server IPs for existing Firewall Policies in Arsys Baremetal
---

# arsys-baremetal\_firewall\_policy\_server_ips

Adds Server IPs to an existing **Firewall Policy** in Arsys baremetal.

## Example Usage

### Add Server IPs to an existing Firewall Policy

```hcl
resource "arsys-baremetal_firewall_policy_server_ips" "example" {
  id = "4C7A872E2C3088E3ADEBB9126489C9E3"
  server_ips = ["08D42710D6071AE597C46EF3C6EB2272", "EC4C4DD06D95333CB63B447D8CA43528"]
}
```

## Argument Reference

The following arguments are supported:

* `id` - (Required) The Id of the existing Firewall Policy.
* `server_ips` - (Required) A list of server_ips to be added to the Firewall Policy. Each element is an id of a server
  IP that you want to assign to the Firewall Policy.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The Id of the Firewall Policy.
* `name` - The name of the Firewall Policy.
* `description` - The description of the Firewall Policy.
* `state` - The current state of the Firewall Policy.
* `creation_date` - The date when the Firewall Policy was created.
* `default` - Define default panel firewalls.
* `cloudpanel_id` - Identifier of the cloud panel.
* `rules` - Rules of the firewall policy.
* `server_ips` - Servers assigned to firewall policy.

## Notes

* This resource cannot be updated. Any changes require replacement.
* The delete operation only removes the resource from the Terraform state.

---

## Import

This resource cannot be imported as it represents an operation to add a server ip to an existing policy.

# Destroy Examples

## Destroy firewall policy server ip remove resource

```shell
terraform destroy -target=arsys-baremetal_firewall_policy_server_ips.example
```