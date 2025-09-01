# Output for all server appliances
output "all_server_appliances" {
  value = data.arsys-baremetal_server_appliances.all
}

# Output for a specific server appliance by Id
output "server_appliance" {
  value = data.arsys-baremetal_server_appliance.example
}