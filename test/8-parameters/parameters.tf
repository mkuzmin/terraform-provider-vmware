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
    parameter1 = "value1"
    parameter2.with.dots = "value2"

    parameter.crash1 = 1
    parameter.crash2 = 1
  }
}
