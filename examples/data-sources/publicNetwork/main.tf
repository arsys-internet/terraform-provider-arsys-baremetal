// Example to discover available public networks
data "arsys-baremetal_public_networks" "all" {}

// Example to retrieve a specific public network by ID
data "arsys-baremetal_public_network" "example" {
  id = var.public_network_id
}