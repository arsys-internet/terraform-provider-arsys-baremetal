# Output for all private networks
output "all_private_networks" {
  value = data.arsys-baremetal_private_networks.all
}

# Output for a specific private network by Id
output "private_network" {
  value = data.arsys-baremetal_private_network.example
}

# Output for a list of all servers assigned to a specific private network
output "all_private_network_servers" {
  value = data.arsys-baremetal_private_network_servers.all
}

# Output for a specific server assigned to a private network
output "private_network_server" {
  value = data.arsys-baremetal_private_network_server.example
}