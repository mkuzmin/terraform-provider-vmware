# Terraform-vSphere Provider

> **Several vSphere providers were developed in parallel. This plugin was first,
> but Terraform team adopted another implementation. Bad luck, yeah.**
>
> **It provides less features, but can be more stable in some cases.

This is a plugin for HashiCorp [Terraform](https://terraform.io/), which helps launching virtual machines on VMware vCenter.

## Usage

- Download the plugin from [Releases](https://github.com/mkuzmin/terraform-vsphere/releases) page.
- [Install](https://terraform.io/docs/plugins/basics.html) it, or put into a directory with configuration files.
- Create a minimal configuration file `web.tf`:
```
provider "vmware" {
  vcenter_server = "vcenter.domain.local"
  user = "root"
  password = "vmware"
}
resource "vmware_virtual_machine" "web" {
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
- `datastore` - by default the same, as base VM.
- `linked_clone` - if *false* (default), full clone is performed from a current state. If *true*, machine is created as a [linked clone](https://pubs.vmware.com/vcd-51/topic/com.vmware.vcloud.admin.doc_51/GUID-4C232B62-4C95-44FF-AD8F-DA2588A5BACC.html) from latest snapshot of base VM.
- `cpus` - a number of CPU sockets in the new VM. By default the same, as base VM.
- `memory` - RAM size in MB. By default the same, as base VM.
- `configuration_parameters` - custom VM parameters.
- `power_on` - if *true* (default), start the newly created machine, and wait till guest OS reports its IP address. VMware Tools must be installed. Timeout is 15 minutes.
- `domain` - enables guest VM customization and hostname renaming. The value specifies domain suffix, hostname is got from `name`.
- `ip_address` - static IP address.
- `subnet_mask` - required if `ip_address` is specified.
- `gateway` - used together with `ip_address`.

## Computed Parameters
- `ip_address` - if `power_on=true` and VMware Tools are installed in guest OS.

## Complete Example
```
provider "vmware" {
  vcenter_server = "vcenter.domain.local"
  user = "domain\user"
  password = "secret"
  insecure_connection = true
}

resource "vmware_virtual_machine" "frontend" {
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

  domain = "vsphere.test"
  ip_address = "192.168.1.10"
  subnet_mask = "255.255.255.0"
  gateway = "192.168.1.1"

  provisioner "remote-exec" {
    connection {
      user = "user"
      password = "secret"
    }
    inline = "uname"
  }

}

output "address" {
  value = "${vmware_virtual_machine.frontend.ip_address}"
}
```
