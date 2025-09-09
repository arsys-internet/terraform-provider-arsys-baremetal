# Example of creating a public IP in the Bare Metal Cloud provider.
resource "arsys-baremetal_public_ip" "example" {
  reverse_dns   = "dns1324.domain.com"
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
  type          = "IPV4"
}

# Example of updating a existing public IP in the Bare Metal Cloud provider.
resource "arsys-baremetal_public_ip" "example" {
  id          = var.public_ip_id
  reverse_dns = "dns5678.domain.com"
}

# Example of importing an existing public IP by ID
resource "arsys-baremetal_public_ip" "test_import" {
}