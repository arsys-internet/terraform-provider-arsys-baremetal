# Output for all baremetal models
output "all_baremetal_models" {
  value = data.arsys-baremetal_baremetal_models.all
}

# Output for a specific baremetal model by ID
output "baremetal_model_by_id" {
  value = data.arsys-baremetal_baremetal_model.example
}