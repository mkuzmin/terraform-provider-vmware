provider "vmware" {
  vcenter_server = "vcenter.vsphere5.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "empty"
  power_on = false
  configuration_parameters = {
    parameter1 = "value"
    parameter2_with_dots = "value"
    parameter_crash1 = 1
    parameter_crash2 = 1
  }
}
