# Public Network

This example demonstrates how to create and manage public networks using the Arsys Baremetal provider.

## What This Example Does

- Creates a public network using `arsys-baremetal_public_network`
- Updates a public network using `arsys-baremetal_public_network`
- Deletes a public network using `terraform destroy -target=arsys-baremetal_public_network.name_of_public_network`
- Import an existing public network using
  `terraform import arsys-baremetal_public_network.name_of_public_network id_of_public_network`

## Other examples associated with the resource

- Associate and disassociate servers to public network using `arsys-baremetal_public_network_servers`
- Associate and disassociate IPs to public network using `arsys-baremetal_public_network_ips`