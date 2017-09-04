provider "vmware" {
  vcenter_server = "vcenter.vsphere55.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "alpine"
  host = "esxi-1.vsphere55.test"
  power_on = false
}
