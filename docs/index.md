---
page_title: "Arsys Baremetal Provider"
description: |-
  The Arsys Baremetal provider is used to manage baremetal infrastructure resources through the Arsys/CloudBuilder API.
---

# Arsys Baremetal Provider

The Arsys Baremetal provider is used to manage baremetal infrastructure resources
(servers, networks, public IPs, firewall policies, SSH keys and more) through the
Arsys/CloudBuilder API.

Use the navigation to the left to read about the available resources and data sources.

## Example Usage

```terraform
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "arsys-internet/arsys-baremetal"
      version = "~> 0.7"
    }
  }
}

provider "arsys-baremetal" {
  # token = "your-api-token"   # prefer the BAREMETAL_API_TOKEN environment variable
}
```

## Authentication

The provider needs an API token to authenticate against the Arsys/CloudBuilder API.
You can provide it in two ways:

- Environment variable (recommended): `export BAREMETAL_API_TOKEN="your-api-token"`
- Statically in the provider block via the `token` argument.

~> **Note:** You must whitelist the source machine IP in your Baremetal panel to
allow API access for your user.

## Schema

### Optional

- `host` (String) Base URL of the Baremetal API. Defaults to
  `https://api.cloudbuilder.es/v1`. Can also be set with the `BAREMETAL_HOST`
  environment variable.
- `token` (String, Sensitive) API token used to authenticate against the Baremetal
  API. Can also be set with the `BAREMETAL_API_TOKEN` environment variable.
