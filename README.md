# Terraform-vSphere Provider

Create a virtual machine on VMware vCenter by cloning an existing VM or template.

## Usage

- Compile and install the plugin
- Create a configuration file:
```
provider "vsphere" {
    server = "vcenter-server"
    user = "account"
# or set VSPHERE_USER environment variable
    password = "secret"
# or set VSPHERE_PASSWORD environment variable
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
## TODO

- customize disk size
- customize RAM size
- Apply customization spec (change hostname)
- get rid of full source path
- refresh for folder, host and pool
- update action
- template = true
- make new snapshot and convert to a template (for future linked clones)
- full clone from current state (option)
