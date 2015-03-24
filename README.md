# Terraform-vSphere Provider

This is a plugin for HashiCorp [Terraform](https://terraform.io/), which helps launching virtual machines on VMware vCenter.

## Usage

- Download the plugin from [Releases](https://github.com/mkuzmin/terraform-vsphere/releases) page.
- [Install](https://terraform.io/docs/plugins/basics.html) it, or put into a directory with configuration files.
- Create a minimal configuration file `web.tf`:
```
provider "vsphere" {
  vcenter_server = "vcenter.domain.local"
  user = "root"
  password = "vmware"
}
resource "vsphere_virtual_machine" "web" {
  name =  "web-1"
  image = "web-base"
}
```
- Modify connection settings, a *name* for the new virtual machine, and a name of a machine to clone from (*image*).
- Run:
```
$ terraform apply
```

## Mandatory Parameters
- `vcenter_server` - address in `hostname[:port]` format.
- `user` - alternatively can be specified via `VSPHERE_USER` environment variable.
- `password` - alternatively can be specified via `VSPHERE_PASSWORD` environment variable.
- `insecure_connection` - Do not check vCenter server SSL certificate. *False* by default.

- `name` of new virtual machine.
- `image` - A name of a base VM or template to clone from. Should include a path in VM folder hierarchy: `folder/subfolder/vm-name`.

## Optional Parameters
- `datacenter` - required, if vSphere has several datacenters.
- `folder` - VM folder to create the virtual machine in `folder/subfolder` format. By default the same, as base VM.
- `host` - by default the same, as base VM. Required if base image is a template.
- `resource_pool` in `pool/nested-pool` format. By default a root of the host.
- `linked_clone` - if *false* (default), full clone is performed from a current state. If *true*, machine is created as a [linked clone](https://pubs.vmware.com/vcd-51/topic/com.vmware.vcloud.admin.doc_51/GUID-4C232B62-4C95-44FF-AD8F-DA2588A5BACC.html) from latest snapshot of base VM.
- `cpus` - a number of CPU sockets in the new VM. By default the same, as base VM.
- `memory` - RAM size in MB. By default the same, as base VM.
- `configuration_parameters` - custom VM parameters.
- `power_on` - if *true* (default), start the newly created machine, and wait till guest OS reports its IP address. VMware Tools must be installed. Timeout is 15 minutes.

## Computed Parameters
- `ip_address` - if `power_on=true` and VMware Tools are installed in guest OS.

## Complete Example
```
provider "vsphere" {
  vcenter_server = "vcenter.domain.local"
  user = "domain\user"
  password = "secret"
  insecure_connection = true
}

resource "vsphere_virtual_machine" "frontend" {
  name =  "frontend-1"
  image = "App/Templates/frontend-base"

  datacenter = "EU"
  folder = "App/Instances"
  host = "esxi1.domain.local"
  resource_pool = "App"

  linked_clone = true
  cpus = 2
  memory = 8192
  configuration_parameters = {
    isolation.tools.copy.disable = "false"
    isolation.tools.paste.disable = "false"
  }

  power_on = true
}

output "address" {
  value = "${vsphere_virtual_machine.frontend.ip_address}"
}
```
