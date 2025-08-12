//Example to create a baremetal server with minimal configuration
resource "arsys-baremetal_server" "basic_server" {
  name          = "test-baremetal-server"
  appliance_id  = var.rocky_linux_appliance_id
  datacenter_id = var.datacenter_spain_id

  hardware = {
    baremetal_model_id = var.baremetal_model_id
  }
}

// Example to create a baremetal server with complete configuration
resource "arsys-baremetal_server" "complete_server" {
  name          = "baremetal-complete-server"
  description   = "Complete baremetal server configuration with Terraform"
  appliance_id  = var.rocky_linux_appliance_id
  datacenter_id = var.datacenter_spain_id

  hardware = {
    baremetal_model_id = var.baremetal_model_id
  }

  password             = "MySecurePassword123!"
  power_on             = true
  install_backup_agent = true

  # Network and security configuration
  firewall_policy_id   = "A1B2C3D4E5F6789012345678901234AB"
  ip_id                = "B2C3D4E5F6789012345678901234ABCD"
  load_balancer_id     = "C3D4E5F6789012345678901234ABCDEF"
  monitoring_policy_id = "D4E5F6789012345678901234ABCDEF12"
  availability_zone_id = "E5F6789012345678901234ABCDEF1234"
}


// Example to update a server
resource "arsys-baremetal_server" "complete_server" {
  name        = "baremetal-complete-server-update"
  description = "Updated baremetal server configuration with Terraform"
}

# Example to import an existing server by ID
resource "arsys-baremetal_server" "server_import" {}
// Execute terraform import command to import the existing server
