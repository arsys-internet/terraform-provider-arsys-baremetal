# Public IP
This example demonstrates how to create and manage public ips using the Arsys Baremetal provider.

## What This Example Does
- Create a public ip using `arsys-baremetal_public_ip`
- Update a specific public ip using `arsys-baremetal_public_ip`
- Delete a specific public ip using `terraform destroy -target=arsys-baremetal_public_ip.name_of_public_ip`
- Import an existing public ip using `terraform import arsys-baremetal_public_ip.name_of_public_ip <id_of_public_ip>`