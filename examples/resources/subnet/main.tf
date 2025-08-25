// Example to create a subnet
resource "arsys-baremetal_subnet" "test" {
  mask          = 28
  datacenter_id = var.datacenter_spain_id
}
