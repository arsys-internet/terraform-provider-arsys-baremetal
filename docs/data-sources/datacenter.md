---
subcategory: "Infrastructure"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Datacenter Data Sources"
sidebar_current: "docs-datasource-datacenter"
description: |-
  Get information on an Arsys Baremetal Datacenter.
---

# arsys-baremetal\_datacenters

The **Datacenter data source** can be used to search and return all existing datacenters.
Also, can be used to search for and return an existing datacenter.
You can provide a string for the id parameter which will be compared with the created datacenter.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all Datacenters

```hcl
data "arsys-baremetal_datacenters" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `datacenters` - List of datacenters

### Get By ID

```hcl
data "arsys-baremetal_datacenter" "example" {
  id = datacenter_id
}
```

## Argument Reference

* `id` - (Required)ID of an existing datacenter that you want to search for.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the Data Center
* `location` - The regional location where the Data Center will be created
* `country_code` - The country code of the Data Center
* `default` - Indicates if the Data Center is the default one in the panel
