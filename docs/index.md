---
layout: "arsys-baremetal"
page_title: "Provider: Arsys Baremetal"
sidebar_current: "docs-arsys-baremetal-index"
description: |-
  The Arsys Baremetal provider is used to manage baremetal infrastructure resources on Arsys.
---

# Arsys Baremetal Provider

The Arsys Baremetal provider is used to interact with the resources supported by
[Arsys Baremetal](https://www.arsys.es/). The provider needs to be configured
with the proper credentials before it can be used.

Use the navigation to the left to read about the available resources and data sources.

## Example Usage

```hcl
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "arsys-internet/arsys-baremetal"
      version = "~> 0.1"
    }
  }
}

# Configure the Arsys Baremetal Provider
provider "arsys-baremetal" {
  token = var.baremetal_api_token
}

```

## Authentication

The provider can be configured with credentials in two ways.

### Static credentials

Credentials can be provided directly in the `provider` block:

```hcl
provider "arsys-baremetal" {
  token = "your-api-token"
}
```

~> **Warning:** Hard-coding credentials into any Terraform configuration is not
recommended, and risks secret leakage should this file ever be committed to a
public version control system. Prefer environment variables or a secrets manager.

### Environment variables

Credentials can also be provided through environment variables, leaving the
`provider` block empty:

```hcl
provider "arsys-baremetal" {}
```

```shell
export BAREMETAL_API_TOKEN="your-api-token"
```

## Schema

### Optional

- `token` (String, Sensitive) API token used to authenticate against the Arsys
  Baremetal API. May also be provided via the `BAREMETAL_API_TOKEN` environment
  variable.