# Example to create a public network
resource "arsys-baremetal_public_network" "example_public_network" {
  public_name   = "p_network_terraform0"
  description   = "Public network for Terraform"
  datacenter_id = "81DEF28500FBC2A973FC0C620DF5B721"
}

# Example to update a public network
resource "arsys-baremetal_public_network" "example_public_network" {
  public_name = "p_network_terraform_update"
  description = "Updated public network for Terraform"
}

# Example to import an existing public network by its ID
resource "arsys-baremetal_public_network" "test_public_network_import" {
}
# Execute terraform import command to import the existing public network



# Example to associate servers to public network
resource "arsys-baremetal_public_network_servers" "example_public_network_server" {
  id = var.public_network_id
  servers = ["08CB1405D8732066D531486E7B7AAC30",
  "08CB1405D8732066D531486E7B7AAC31"]
}

# Example to disassociate alls server from public network
resource "arsys-baremetal_public_network_servers" "example_public_network_server" {
  id      = var.public_network_id
  servers = []
}

# Example to disassociate one server (08CB1405D8732066D531486E7B7AAC31) from public network
resource "arsys-baremetal_public_network_servers" "example_public_network_server" {
  id      = var.public_network_id
  servers = ["08CB1405D8732066D531486E7B7AAC30"]
}



# Example to associate IPs to public network
resource "arsys-baremetal_public_network_ips" "example_public_network_ip" {
  public_network_id = var.public_network_id
  action            = true
  ips = ["08CB1405D8732066D531486E7B7AAC30",
  "08CB1405D8732066D531486E7B7AAC31"]
}

# Example to disassociate IPs from public network
resource "arsys-baremetal_public_network_ips" "example_public_network_ip" {
  public_network_id = var.public_network_id
  action            = false
  ips = ["08CB1405D8732066D531486E7B7AAC30",
  "08CB1405D8732066D531486E7B7AAC31"]
}

