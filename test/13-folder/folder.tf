provider "vmware" {
  vcenter_server = "vcenter.vsphere55.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}

resource "vmware_vm_folder" "test" {
  datacenter = "dc1"
  parent =  "/folder1"
  name =  "test"
}
