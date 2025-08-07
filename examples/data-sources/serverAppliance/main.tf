# Example of using the server appliance data source in the Bare Metal Cloud provider.
data "arsys-baremetal_server_appliances" "all" {}

# Example of using the server appliance data source in the Bare Metal Cloud provider with a specific ID.
data "arsys-baremetal_server_appliance" "example" {
  id = var.server_appliance_id
}