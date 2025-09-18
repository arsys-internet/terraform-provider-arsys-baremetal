---
subcategory: "Compute Engine"
layout: "arsys-baremetal"
page_title: "Arsys Baremetal : Server appliance Data source"
sidebar_current: "docs-datasource-server_appliance"
description: |-
  Get information about Server Appliances in Arsys Baremetal .
---

# arsys-baremetal\_server\_appliance

The **Server appliances data source** can be used to search and return all existing server appliances.
Also, can be used to search for and return an existing server appliance by id.
You can provide a string for the id parameter which will be compared with created server appliances.
If a single match is found, it will be returned. If it is not found an error will be returned.

## Example Usage

### Get all Server Appliances

```hcl
data "arsys-baremetal_server_appliances" "all" {}
```

## Attributes Reference

The following attributes are returned by the datasource:

* `server_appliances` - List of server appliances

### Get By ID

```hcl
data "arsys-baremetal_server_appliance" "example" {
  id = server_appliance_id
}
```

## Argument Reference

* `id` - (Required) ID of an existing server appliance that you want to search for.
  `id` must be provided. If none, the datasource will return an error.

## Attributes Reference

The following attributes are returned by the datasource:

* `id` - Identifier of the server appliance
* `name` - The name of the server appliance
* `available_datacenters` - List of available datacenters
* `os_family` - Family of the operating system
* `os` - The operating system of the server appliance
* `os_version` - Version of the operating system
* `os_architecture` - The architecture of the operating system (e.g., "x64", "arm64")
* `os_image_type` - Image installation type (e. g., "STANDARD", "MINIMAL", "ISO_OS", "ISO_TOOL")
* `type` - The type of the server appliance (e.g., "IMAGE", "APPLICATION", "ISO")
* `server_type_compatibility` - List of servers type compatible with the image
* `min_hdd_size` - Minimum required hdd size for this image
* `licenses` - List of image licenses acquired
* `version` - Version of the application
* `categories` - List of categories to which the application belongs