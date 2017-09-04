provider "vmware" {
  vcenter_server = "vcenter.vsphere55.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}
resource "vmware_virtual_machine" "vm" {
  name =  "vm-1"
  image = "basic"
  linked_clone = true
  provisioner "remote-exec" {
    connection {
      user = "jetbrains"
      password = "jetbrains"
      agent = false
    }
    inline = "uname -a"
  }
}
