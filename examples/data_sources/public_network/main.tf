# Example to discover available public networks
data "arsys-baremetal_public_networks" "all" {}

# Example to retrieve a specific public network by ID
data "arsys-baremetal_public_network" "example" {
  id = var.public_network_id
}

# Example to retrieve all Ips in a specific public network
data "arsys-baremetal_public_network_ips" "all" {
  public_network_id = var.public_network_id
}

# Example to retrieve a specific IP by Id in a specific public network
data "arsys-baremetal_public_network_ip" "example" {
  public_network_id = var.public_network_id
  id = var.ip_id
}