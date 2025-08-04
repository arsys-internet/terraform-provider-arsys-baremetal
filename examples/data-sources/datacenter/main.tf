//Data source to discover available datacenters to create your resources
data "arsys-baremetal_datacenters" "all" {}

// Data source to retrieve a specific datacenter by ID
data "arsys-baremetal_datacenter" "example" {
  id = var.datacenter_spain_id
}