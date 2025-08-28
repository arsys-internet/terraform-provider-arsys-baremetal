# Example of using the subnets data source in the Bare Metal Cloud provider.
data "arsys-baremetal_subnets" "all" {}

# Example of using the subnet data source in the Bare Metal Cloud provider with a specific Id.
data "arsys-baremetal_subnet" "example" {
  id = var.subnet_id
}
