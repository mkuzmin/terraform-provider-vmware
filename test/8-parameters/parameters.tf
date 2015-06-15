provider "vsphere" {
  vcenter_server = "vcenter.vsphere.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}
resource "vsphere_virtual_machine" "vm" {
  name =  "vm-1"
  image = "empty"
  power_on = false
  configuration_parameters = {
    parameter1 = "value"
    parameter2.with.dots = "value"
    parameter.crash1 = 1
    parameter.crash2 = 1
  }
}
