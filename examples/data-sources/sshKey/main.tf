# Example of using the SSH key data source in the Bare Metal Cloud provider.
data "arsys-baremetal_ssh_keys" "all" {}

# Example of using the SSH key data source in the Bare Metal Cloud provider with a specific ID.
data "arsys-baremetal_ssh_key" "example" {
  id = var.ssh_key_id
}