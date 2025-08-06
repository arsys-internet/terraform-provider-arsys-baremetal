// Example to discover available private networks
data "arsys-baremetal_private_networks" "all" {}

// Example to retrieve a specific private network by ID
data "arsys-baremetal_private_network" "example" {
  id = "private_netowrk_id"
}