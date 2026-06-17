# Example of a Terraform configuration using the arsys-baremetal provider
terraform {
  required_providers {
    arsys-baremetal = {
      source  = "arsys-internet/arsys-baremetal"
      version = "~> 0.1"
    }
  }
}

# Configure the provider.
# Credentials can be set in the provider block or via the BAREMETAL_HOST and
# BAREMETAL_API_TOKEN environment variables.
provider "arsys-baremetal" {
  # host  = "https://api.cloudbuilder.es/v1"
  # token = "your-api-token"
}