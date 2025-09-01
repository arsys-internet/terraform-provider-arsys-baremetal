# Output for all SSH keys
output "all_ssh_keys" {
  value = data.arsys-baremetal_ssh_keys.all
}

# Output for a specific SSH key by Id
output "ssh_key" {
  value = data.arsys-baremetal_ssh_key.example
}