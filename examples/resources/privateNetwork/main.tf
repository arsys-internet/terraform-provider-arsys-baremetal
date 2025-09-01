# Example to create a private network
resource "arsys-baremetal_private_network" "example_private_network" {
  name            = "p_network_terraform0"
  network_address = "192.168.4.0"
  subnet_mask     = "255.255.255.0"
  datacenter_id   = "81DEF28500FBC2A973FC0C620DF5B721"
  description     = "Private network for Terraform"
}

# Example to update a private network
resource "arsys-baremetal_private_network" "example_private_network" {
  name            = "p_network_terraform_update"
  description     = "Updated private network for Terraform"
  network_address = "192.168.5.0"
  subnet_mask     = "255.255.255.0"
}

# Example to import an existing private network by its ID
resource "arsys-baremetal_private_network" "test_private_network_import" {
}
# Execute terraform import command to import the existing private network