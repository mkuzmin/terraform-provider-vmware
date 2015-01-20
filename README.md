# Terraform-vSphere Provider

Create a virtual machine on VMware vCenter by cloning an existing VM or template.

## Usage

- Compile and install the plugin
- Create a configuration file:
```
provider "vsphere" {
    server = "vcenter-server"
    user = "account"
// or set VSPHERE_USER environment variable
    password = "secret"
// or set VSPHERE_PASSWORD environment variable
}

resource "vsphere_vm" "machine" {
    name =  "machine1"
    source = "Full/Path/to/VM"
    datacenter = "DC"
    folder = "Full/Path/to/Folder"
    host = "hostname"
    pool = "Resource/Pool"
// optional
    linked_clone = true
}
```
- Run

    $ terraform apply

## TODO

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
- full clone from snapshot or current state
