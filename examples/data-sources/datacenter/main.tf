# Example to discover available datacenters to create your resources
data "arsys-baremetal_datacenters" "all" {}

# Example to retrieve a specific datacenter by Id
data "arsys-baremetal_datacenter" "example" {
  id = var.datacenter_spain_id
}