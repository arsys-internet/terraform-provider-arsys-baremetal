# Output for a specific private network by ID
output "private_network" {
  value = arsys-baremetal_private_network.example_private_network
}

# Output for a private network servers assign
output "private_network_servers_assign" {
  value = arsys-baremetal_private_network_servers_assign.example_assign
}

# Output for a private network server remove
output "private_network_server_remove" {
  value = arsys-baremetal_private_network_server_remove.example_remove
}
