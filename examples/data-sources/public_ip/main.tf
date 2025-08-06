# Example of using the public IP data source in the Bare Metal Cloud provider.
data "arsys-baremetal_public_ips" "all" {}

# Example of using the public IP data source in the Bare Metal Cloud provider with a specific ID.
data "arsys-baremetal_public_ip" "example" {
  id = var.public_ip_id
}