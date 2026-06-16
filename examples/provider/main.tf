# Example of a Terraform configuration using the arsys-baremetal provider
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "arsys-internet/arsys-baremetal"
      version = "~> 0.7"
    }
  }
}

# Set your API token
# You can set the token directly in the provider block or use the
# BAREMETAL_API_TOKEN (and optionally BAREMETAL_HOST) environment variables.
provider "arsys-baremetal" {
  # token = "token"
}
