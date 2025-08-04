// Example of a Terraform configuration using the arsys-baremetal provider
terraform {
  required_providers {
    arsys-baremetal = {
      source = "local/arsys-baremetal"
    }
  }
}

# Set your API host and token
# You can set the host and token directly in the provider block or use environment variables ARSYS_BAREMETAL_HOST and ARSYS_BAREMETAL_TOKEN.
provider "arsys-baremetal" {
  #   host     = "https://api.ejemplo.com"
  # token = "token"
}
