# Output for all private networks
output "all_private_networks" {
  value = data.arsys-baremetal_private_networks.all
}

# Output for a specific private network by ID
output "private_network" {
  value = data.arsys-baremetal_private_network.example
}