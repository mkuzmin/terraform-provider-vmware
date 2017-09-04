provider "vmware" {
  vcenter_server = "vcenter.vsphere55.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "empty"
  host = "esxi-4.vsphere55.test"
  datastore = "datastore4-2"
  power_on = false
}
