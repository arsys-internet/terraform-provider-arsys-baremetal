# Output for all subnets
output "all_subnets" {
  value = data.arsys-baremetal_subnets.all
}

# Output for a specific subnet by ID
output "subnet" {
  value = data.arsys-baremetal_subnet.example
}