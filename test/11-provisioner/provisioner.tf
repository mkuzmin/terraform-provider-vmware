provider "vsphere" {
  vcenter_server = "vcenter.vsphere.test"
  user = "root"
  password = "jetbrains"
  insecure_connection = true
}
resource "vsphere_virtual_machine" "vm" {
  name =  "vm-1"
  image = "basic"
  linked_clone = true
  provisioner "remote-exec" {
    connection {
      user = "jetbrains"
      password = "jetbrains"
    }
    inline = "uname -a"
  }
}
