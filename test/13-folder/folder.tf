provider "vmware" {
  vcenter_server = "vcenter.vsphere5.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}

resource "vmware_vm_folder" "test" {
  datacenter = "DC1"
  parent =  "folder1"
  name =  "test"
}
