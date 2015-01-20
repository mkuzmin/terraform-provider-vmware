# Terraform-vSphere Provider

Create a virtual machine on VMware vCenter by cloning an existing VM or template.

## Usage

- Compile and install the plugin
- Create a configuration file:
```
provider "vsphere" {
    server = "vcenter-server"
    user = "account"
    password = "secret"
}

resource "vsphere_vm" "machine" {
    name =  "machine1"
    source = "Full/Path/to/VM"
    datacenter = "DC"
    folder = "Full/Path/to/Folder"
    host = "hostname"
    pool = "Resource/Pool"
}
```
- Run
    $ terraform apply

## Environment Variables

Instead of storing credentials in the file, `VSPHERE_USER` and `VSPHERE_PASSWORD` environment variables can be used.

## TODO

- linked_clone = true
- power_on = false
- template = true
- customize disk size
- customize RAM size
- Apply customization spec (change hostname)
- get rid of full source path
- persist UUID instead of vm name
- read/update/delete actions
- partitial update
- move datacenter name into provider settings
- making new snapshots
