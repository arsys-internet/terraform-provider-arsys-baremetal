# Example to discover baremetal models to create your resources
data "arsys-baremetal_baremetal_models" "all" {}

# Example to retrieve a specific baremetal model by ID
data "arsys-baremetal_baremetal_model" "example" {
  id = var.bmc_ssd_id
}