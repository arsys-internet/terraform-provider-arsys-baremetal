// Example to discover available baremetal servers
data "arsys-baremetal_servers" "all" {}

// Example to retrieve a specific baremetal server by ID
data "arsys-baremetal_server" "server" {
  id = "A982AE3D56CEB4CEB2FB9B62C6B74691"
}