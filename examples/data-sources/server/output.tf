# Output for all baremetal servers
output "all_servers" {
  value = data.arsys-baremetal_servers.all
}

# Output for a specific server by Id
output "server" {
  value = data.arsys-baremetal_server.server
}