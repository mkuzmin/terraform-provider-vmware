package main

import (
	"fmt"
	"log"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

func resourceVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineCreate,
		Read:   resourceVirtualMachineRead,
		Delete: resourceVirtualMachineDelete,

		Schema: map[string]*schema.Schema{
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"source": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"folder": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
                Computed: true,
			},
			"pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
                Computed: true,
            },
			"linked_clone": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
			},
			"cpus": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
                Computed: true,
			},
			"memory": &schema.Schema{
				Type:     schema.TypeInt,
				Optional: true,
                Computed: true,
			},
			"power_on": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
                Default:  true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*govmomi.Client)

	source_vm_ref, err := client.SearchIndex().FindByInventoryPath(fmt.Sprintf("%s/vm/%s", d.Get("datacenter").(string), d.Get("source").(string)))
	if err != nil {
		return fmt.Errorf("Error reading vm: %s", err)
	}
    if source_vm_ref == nil {
        return fmt.Errorf("Cannot find source vm %s", d.Get("source").(string))
    }
    source_vm := source_vm_ref.(*govmomi.VirtualMachine)

	var folder_ref govmomi.Reference
    var source_vm_mo mo.VirtualMachine
    var folder *govmomi.Folder
    if d.Get("folder").(string) != "" {
        folder_ref, err = client.SearchIndex().FindByInventoryPath(fmt.Sprintf("%v/vm/%v", d.Get("datacenter").(string), d.Get("folder").(string)))
        if err != nil {
            return fmt.Errorf("Error reading folder: %s", err)
        }
        if folder_ref == nil {
            return fmt.Errorf("Cannot find folder %s", d.Get("folder").(string))
        }

        folder = folder_ref.(*govmomi.Folder)
    } else {
        err = client.Properties(source_vm.Reference(), []string{"parent"}, &source_vm_mo)
        if err != nil {
            return fmt.Errorf("Error reading parent VM folder")
        }
        folder = govmomi.NewFolder(client, *source_vm_mo.Parent)
    }


	var relocateSpec types.VirtualMachineRelocateSpec

    var pool_mor types.ManagedObjectReference
    if d.Get("host").(string) != "" && d.Get("pool").(string) != ""	{
        pool_ref, err := client.SearchIndex().FindByInventoryPath(fmt.Sprintf("%v/host/%v/Resources/%v", d.Get("datacenter").(string), d.Get("host").(string), d.Get("pool").(string)))
        if err != nil {
            return fmt.Errorf("Error reading resource pool: %s", err)
        }
        if pool_ref == nil {
            return fmt.Errorf("Cannot find resource pool %s", d.Get("pool").(string))
        }
        pool_mor = pool_ref.Reference()
        relocateSpec.Pool = &pool_mor
    }

	if d.Get("linked_clone").(bool) {
		relocateSpec.DiskMoveType = "createNewChildDiskBacking"
	}
    var confSpec types.VirtualMachineConfigSpec
    if d.Get("cpus") != nil {
        confSpec.NumCPUs = d.Get("cpus").(int)
    }
    if d.Get("memory") != nil {
        confSpec.MemoryMB = int64(d.Get("memory").(int))
    }

	cloneSpec := types.VirtualMachineCloneSpec{
		Location: relocateSpec,
        Config:   &confSpec,
        PowerOn:  d.Get("power_on").(bool),
	}
    if d.Get("linked_clone").(bool) {
        err = client.Properties(source_vm.Reference(), []string{"snapshot"}, &source_vm_mo)
        if err != nil {
            return fmt.Errorf("Error reading snapshot")
        }
        if source_vm_mo.Snapshot == nil {
            return fmt.Errorf("Base VM has no snapshots")
        }
        cloneSpec.Snapshot = source_vm_mo.Snapshot.CurrentSnapshot
    }

	task, err := source_vm.Clone(folder, d.Get("name").(string), cloneSpec)
	if err != nil {
		return fmt.Errorf("Error clonning vm: %s", err)
	}
	info, err := task.WaitForResult(nil)
	if err != nil {
		return fmt.Errorf("Error clonning vm: %s", err)
	}

	vm_mor := info.Result.(types.ManagedObjectReference)
    d.SetId(vm_mor.Value)
    vm := govmomi.NewVirtualMachine(client, vm_mor)
    // workaround for https://github.com/vmware/govmomi/issues/218
    if d.Get("power_on").(bool) {
        ip, err := vm.WaitForIP()
        if err != nil {
            log.Printf("[ERROR] Cannot read ip address: %s", err)
        } else {
            d.Set("ip_address", ip)
        }
    }

    return nil
}

func resourceVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*govmomi.Client)
    vm_mor := types.ManagedObjectReference{Type: "VirtualMachine", Value: d.Id() }
    vm := govmomi.NewVirtualMachine(client, vm_mor)

    var vm_mo mo.VirtualMachine
    err := client.Properties(vm.Reference(), []string{"summary"}, &vm_mo)
    if err != nil {
        d.SetId("")
        return nil
    }
    d.Set("name", vm_mo.Summary.Config.Name)
    d.Set("cpus", vm_mo.Summary.Config.NumCpu)
    d.Set("memory", vm_mo.Summary.Config.MemorySizeMB)

    if vm_mo.Summary.Runtime.PowerState == "poweredOn" {
        d.Set("power_on", true)
    } else {
        d.Set("power_on", false)
    }

    if d.Get("power_on").(bool) {
        ip, err := vm.WaitForIP()
        if err != nil {
            log.Printf("[ERROR] Cannot read ip address: %s", err)
        } else {
            d.Set("ip_address", ip)
        }
    }

	return nil
}

func resourceVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*govmomi.Client)
    vm_mor := types.ManagedObjectReference{Type: "VirtualMachine", Value: d.Id() }
    vm := govmomi.NewVirtualMachine(client, vm_mor)

    task, err := vm.PowerOff()
    if err != nil {
        return fmt.Errorf("Error powering vm off: %s", err)
    }
    task.WaitForResult(nil)

    task, err = vm.Destroy()
    if err != nil {
        return fmt.Errorf("Error deleting vm: %s", err)
    }
    _, err = task.WaitForResult(nil)
    if err != nil {
        return fmt.Errorf("Error deleting vm: %s", err)
    }

    return nil
}
