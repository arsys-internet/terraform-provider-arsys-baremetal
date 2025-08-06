// Example to discover available private networks
data "arsys-baremetal_private_networks" "all" {}

// Example to retrieve a specific private network by ID
data "arsys-baremetal_private_network" "example" {
  id = "192117615D6F725215A21B05C87068BB"
}