provider "vmware" {
  vcenter_server = "vcenter.vsphere5.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}

resource "vmware_vm_folder" "test" {
  parent =  "/DC1/vm/folder1"
  name =  "test"
}
