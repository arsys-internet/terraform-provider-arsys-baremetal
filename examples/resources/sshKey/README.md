# SSH key

This example demonstrates how to create and manage SSH keys using the Arsys Baremetal provider.

## What This Example Does

- Create a SSH key using `arsys-baremetal_ssh_key`
- Update a specific SSH key using `arsys-baremetal_ssh_key`
- Delete a specific SSH key using `terraform destroy -target=arsys-baremetal_ssh_key.name_of_ssh_key`
- Import an existing SSH key using `terraform import arsys-baremetal_ssh_key.name_of_ssh_key <id_of_ssh_key>`