# Output for all datacenters
output "all_datacenters" {
  value = data.arsys-baremetal_datacenters.all
}

# Output for a specific datacenter by ID
output "datacenter" {
  value = data.arsys-baremetal_datacenter.example
}