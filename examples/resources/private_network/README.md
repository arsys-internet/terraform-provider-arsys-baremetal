# Private Network

This example demonstrates how to create and manage private networks using the Arsys Baremetal provider.

## What This Example Does

- Creates a private network using `arsys-baremetal_private_network`
- Updates a private network using `arsys-baremetal_private_network`
- Deletes a private network using `terraform destroy -target=arsys-baremetal_private_network.name_of_private_network`
- Import an existing private network using `terraform import arsys-baremetal_private_network.name_of_private_network id_of_private_network`
- Assign servers to a private network using `arsys-baremetal_private_network_servers_assign`
- Remove server from a private network using `arsys-baremetal_private_network_server_remove`