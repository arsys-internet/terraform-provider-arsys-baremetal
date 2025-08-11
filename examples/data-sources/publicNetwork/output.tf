//Output for all public networks
output "all_public_networks" {
  value = data.arsys-baremetal_public_networks.all
}

// Output for a specific public network by its ID
output "public_network" {
  value = data.arsys-baremetal_public_network.example
}