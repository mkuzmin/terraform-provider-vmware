# Terraform-vSphere Provider

Create a virtual machine on VMware vCenter by cloning an existing VM or template.

## Usage

- Download the plugin from [Releases](https://github.com/mkuzmin/terraform-vsphere/releases) page
- [Install](https://terraform.io/docs/plugins/basics.html) it, or put to the same directory with configuration files.
- Create a configuration file:
```
provider "vsphere" {
    server = "vcenter-server"
# or set VSPHERE_USER environment variable
    user = "account"
# or set VSPHERE_PASSWORD environment variable
    password = "secret"
}

resource "vsphere_virtual_machine" "machine" {
    name =  "machine1"
    source = "Full/Path/to/VM"
    datacenter = "DC"
    folder = "Full/Path/to/Folder"
    host = "hostname"
    pool = "Resource/Pool"
# optional
    linked_clone = true
    # power_on = false
}
```
- Run
```
$ terraform apply

$ terraform destroy
```
