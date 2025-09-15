---
subcategory: "Infrastructure"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal: Baremetal Models Data Sources"
sidebar_current: "docs-datasource_baremetal_models"
description: |-
  Get information about Baremetal Models in Arsys Baremetal.
---

# arsys-baremetal\_baremetal\_models

The **Baremetal models data source** can be used to search and return all existing baremetal models.
Also, it can be used to search for and return an existing baremetal model.
You can provide a string for the id parameter which will be compared with an existing baremetal model.
If a single match is found, it will be returned. If it is not found, an error will be returned.

## Example Usage

### Get all Baremetal Models

```hcl
data "arsys-baremetal_baremetal_models" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `baremetal_models` - List of baremetal models.

### Get By ID

```hcl
data "arsys-baremetal_baremetal_model" "example" {
  id = baremetal_model_id
}
```

## Argument Reference

* `id` - (Required) ID of an existing baremetal model that you want to search for.

`id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the Baremetal Model
* `name` - The name of the Baremetal Model
* `state_id` - Identifier representing the current state of the Baremetal Model
* `state` - Current state of the Baremetal Model (e.g., "ENABLED", "DISABLED")
* `hardware` attribute contains the following nested attributes:
    * `core` - Number of cores
    * `cores_per_processor` - Number of cores per processor
    * `ram` - Amount of RAM
    * `unit` - Unit of measurement for RAM
    * `hdds` - List of hard disk drive configurations
        * `size` - HDD size
        * `unit` - Unit of measurement for the HDD size
        * `disk_type` - Type of disk (e.g., "HDD", "SSD")
        * `disk_raid` - RAID configuration (e.g., "Hardware RAID 1", "Software RAID 5")
        * `disk_raid_count` - Number of disks in RAID configuration
        * `is_main` - Whether this is the main HDD
* `availability` - List of availability information per datacenter, including datacenter ID, availability status,
  connection speeds, and redundancy options
    * `datacenter_id` - Identifier of the datacenter where the baremetal model is available
    * `available` - Whether the baremetal model is available in this datacenter
    * `available_connections_speeds` - List of available connection speeds
    * `available_with_redundancy` - Whether the model is available with redundancy in this datacenter
    * `available_without_redundancy` - Whether the model is available without redundancy in this datacenter