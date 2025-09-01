# Example of a Terraform configuration using the arsys-baremetal provider
terraform {
  required_providers {
    arsys-baremetal = {
      source = "local/arsys-baremetal"
    }
  }
}

# Set your API token
# You can set the token directly in the provider block or use environment variables ARSYS_BAREMETAL_TOKEN.
provider "arsys-baremetal" {
  # token = "token"
}
