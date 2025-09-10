# Example of creating a SSH key in the Bare Metal Cloud provider.
resource "arsys-baremetal_ssh_key" "example" {
  name        = "my-ssh-key"
  description = "SSH key for accessing servers"
  public_key  = "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC7vbqajDhLpOF6JVqHGBSDuBv3wPsX1VQCnXFZPVNhGJ4q test@example.com"
}

# Example of updating an existing SSH key in the Bare Metal Cloud provider.
resource "arsys-baremetal_ssh_key" "example" {
  id          = var.ssh_key_id
  name        = "new my-ssh-key"
  description = "new SSH key for accessing servers"
}

# Example of importing an existing SSH key by ID
resource "arsys-baremetal_ssh_key" "test_import" {
}