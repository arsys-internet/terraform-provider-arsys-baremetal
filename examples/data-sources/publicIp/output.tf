//Output for all public ips
output "all_public_ips" {
  value = data.arsys-baremetal_public_ips.all
}

// Output for a specific public ip by ID
output "public_ip" {
  value = data.arsys-baremetal_public_ip.example
}