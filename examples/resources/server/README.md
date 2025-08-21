# Private Network

This example demonstrates how to create and manage servers using the Arsys Baremetal provider.

## What This Example Does

- Creates a server using `arsys-baremetal_server`
- Updates a server using `arsys-baremetal_server`
- Deletes a server using `terraform destroy -target=arsys-baremetal_server.name_of_server`
- Import an existing server using `terraform import arsys-baremetal_server.name_of_server id_of_server`