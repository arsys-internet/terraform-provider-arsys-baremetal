# Example to discover available private networks
data "arsys-baremetal_private_networks" "all" {}

# Example to retrieve a specific private network by Id
data "arsys-baremetal_private_network" "example" {
  id = var.private_network_id
}

# Example to discover all servers that are assigned to a specific private network
data "arsys-baremetal_private_network_servers" "all" {
  id = var.private_network_id
}

# Example to retrieve a specific server assigned to a private network
data "arsys-baremetal_private_network_server" "example" {
  private_network_id = var.private_network_id
  id                 = var.id
}