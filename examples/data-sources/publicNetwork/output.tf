# Output for all public networks
output "all_public_networks" {
  value = data.arsys-baremetal_public_networks.all
}

# Output for a specific public network by ID
output "public_network" {
  value = data.arsys-baremetal_public_network.example
}

# Output for all IPs in a public network
output "all_ips_public_network" {
  value = data.arsys-baremetal_public_network_ips.all
}

# Output for a specific IP in a public network by ID
output "ip_public_network" {
  value = data.arsys-baremetal_public_network_ip.example
}