package main

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
	"log"
	"strings"
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
				ForceNew: true,
			},
			"image": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"datacenter": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"resource_pool": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datastore": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"linked_clone": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
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
			"domain": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_mask": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gateway": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"configuration_parameters": &schema.Schema{
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"power_on": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	finder := find.NewFinder(client, false)
	ctx := context.TODO()

	dc_name := d.Get("datacenter").(string)
	if dc_name == "" {
		dc, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			return fmt.Errorf("Error reading default datacenter: %s", err)
		}
		var dc_mo mo.Datacenter
		err = dc.Properties(ctx, dc.Reference(), []string{"name"}, &dc_mo)
		if err != nil {
			return fmt.Errorf("Error reading datacenter name: %s", err)
		}
		dc_name = dc_mo.Name
		finder.SetDatacenter(dc)
		d.Set("datacenter", dc_name)
	}

	image_name := d.Get("image").(string)
	image_ref, err := object.NewSearchIndex(client).FindByInventoryPath(ctx, fmt.Sprintf("%s/vm/%s", dc_name, image_name))
	if err != nil {
		return fmt.Errorf("Error reading vm: %s", err)
	}
	if image_ref == nil {
		return fmt.Errorf("Cannot find image %s", image_name)
	}
	image := image_ref.(*object.VirtualMachine)

	var image_mo mo.VirtualMachine
	err = image.Properties(ctx, image.Reference(), []string{"parent", "config.template", "resourcePool", "snapshot", "guest.toolsVersionStatus2", "config.guestFullName"}, &image_mo)
	if err != nil {
		return fmt.Errorf("Error reading base VM properties: %s", err)
	}

	var folder_ref object.Reference
	var folder *object.Folder
	if d.Get("folder").(string) != "" {
		folder_ref, err = object.NewSearchIndex(client).FindByInventoryPath(ctx, fmt.Sprintf("%v/vm/%v", dc_name, d.Get("folder").(string)))
		if err != nil {
			return fmt.Errorf("Error reading folder: %s", err)
		}
		if folder_ref == nil {
			return fmt.Errorf("Cannot find folder %s", d.Get("folder").(string))
		}

		folder = folder_ref.(*object.Folder)
	} else {
		folder = object.NewFolder(client, *image_mo.Parent)
	}

	host_name := d.Get("host").(string)
	if host_name == "" {
		if image_mo.Config.Template == true {
			return fmt.Errorf("Image is a template, 'host' is a required")
		} else {
			var pool_mo mo.ResourcePool
			err = property.DefaultCollector(client).RetrieveOne(ctx, *image_mo.ResourcePool, []string{"owner"}, &pool_mo)
			if err != nil {
				return fmt.Errorf("Error reading resource pool of base VM: %s", err)
			}

			if strings.Contains(pool_mo.Owner.Value, "domain-s") {
				var host_mo mo.ComputeResource
				err = property.DefaultCollector(client).RetrieveOne(ctx, pool_mo.Owner, []string{"name"}, &host_mo)
				if err != nil {
					return fmt.Errorf("Error reading host of base VM: %s", err)
				}
				host_name = host_mo.Name
			} else if strings.Contains(pool_mo.Owner.Value, "domain-c") {
				var cluster_mo mo.ClusterComputeResource
				err = property.DefaultCollector(client).RetrieveOne(ctx, pool_mo.Owner, []string{"name"}, &cluster_mo)
				if err != nil {
					return fmt.Errorf("Error reading cluster of base VM: %s", err)
				}
				host_name = cluster_mo.Name
			} else {
				return fmt.Errorf("Unknown compute resource format of base VM: %s", pool_mo.Owner.Value)
			}
		}
	}

	pool_name := d.Get("resource_pool").(string)
	pool_ref, err := object.NewSearchIndex(client).FindByInventoryPath(ctx, fmt.Sprintf("%v/host/%v/Resources/%v", dc_name, host_name, pool_name))
	if err != nil {
		return fmt.Errorf("Error reading resource pool: %s", err)
	}
	if pool_ref == nil {
		return fmt.Errorf("Cannot find resource pool %s", pool_name)
	}

	var relocateSpec types.VirtualMachineRelocateSpec
	var pool_mor types.ManagedObjectReference
	pool_mor = pool_ref.Reference()
	relocateSpec.Pool = &pool_mor

	datastore_name := d.Get("datastore").(string)
	if datastore_name != "" {
		datastore_ref, err := finder.Datastore(ctx, fmt.Sprintf("/%v/datastore/%v", dc_name, datastore_name))
		if err != nil {
			return fmt.Errorf("Cannot find datastore '%s'", datastore_name)
		}
		datastore_mor := datastore_ref.Reference()
		relocateSpec.Datastore = &datastore_mor
	}

	if d.Get("linked_clone").(bool) {
		relocateSpec.DiskMoveType = "createNewChildDiskBacking"
	}
	var confSpec types.VirtualMachineConfigSpec
	if d.Get("cpus") != nil {
		confSpec.NumCPUs = int32(d.Get("cpus").(int))
	}
	if d.Get("memory") != nil {
		confSpec.MemoryMB = int64(d.Get("memory").(int))
	}

	params := d.Get("configuration_parameters").(map[string]interface{})
	var ov []types.BaseOptionValue
	if len(params) > 0 {
		for k, v := range params {
			key := strings.Replace(k, "_", ".", -1)
			value := v
			o := types.OptionValue{
				Key:   key,
				Value: &value,
			}
			ov = append(ov, &o)
		}
		confSpec.ExtraConfig = ov
	}

	cloneSpec := types.VirtualMachineCloneSpec{
		Location: relocateSpec,
		Config:   &confSpec,
		PowerOn:  d.Get("power_on").(bool),
	}
	if d.Get("linked_clone").(bool) {
		if image_mo.Snapshot == nil {
			return fmt.Errorf("`linked_clone=true`, but image VM has no snapshots")
		}
		cloneSpec.Snapshot = image_mo.Snapshot.CurrentSnapshot
	}

	domain := d.Get("domain").(string)
	ip_address := d.Get("ip_address").(string)
	if domain != "" {
		if image_mo.Guest.ToolsVersionStatus2 == "guestToolsNotInstalled" {
			return fmt.Errorf("VMware tools are not installed in base VM")
		}
		if !strings.Contains(image_mo.Config.GuestFullName, "Linux") && !strings.Contains(image_mo.Config.GuestFullName, "CentOS") {
			return fmt.Errorf("Guest customization is supported only for Linux. Base image OS is: %s", image_mo.Config.GuestFullName)
		}
		customizationSpec := types.CustomizationSpec{
			GlobalIPSettings: types.CustomizationGlobalIPSettings{},
			Identity: &types.CustomizationLinuxPrep{
				HostName: &types.CustomizationVirtualMachineName{},
				Domain:   domain,
			},
			NicSettingMap: []types.CustomizationAdapterMapping{
				{
					Adapter: types.CustomizationIPSettings{},
				},
			},
		}
		if ip_address != "" {
			mask := d.Get("subnet_mask").(string)
			if mask == "" {
				return fmt.Errorf("'subnet_mask' must be set, if static 'ip_address' is specified")
			}
			customizationSpec.NicSettingMap[0].Adapter.Ip = &types.CustomizationFixedIp{
				IpAddress: ip_address,
			}
			customizationSpec.NicSettingMap[0].Adapter.SubnetMask = d.Get("subnet_mask").(string)
			gateway := d.Get("gateway").(string)
			if gateway != "" {
				customizationSpec.NicSettingMap[0].Adapter.Gateway = []string{gateway}
			}
		} else {
			customizationSpec.NicSettingMap[0].Adapter.Ip = &types.CustomizationDhcpIpGenerator{}
		}
		cloneSpec.Customization = &customizationSpec
	} else if ip_address != "" {
		return fmt.Errorf("'domain' must be set, if static 'ip_address' is specified")
	}

	task, err := image.Clone(ctx, folder, d.Get("name").(string), cloneSpec)
	if err != nil {
		return fmt.Errorf("Error clonning vm: %s", err)
	}
	info, err := task.WaitForResult(ctx, nil)
	if err != nil {
		return fmt.Errorf("Error clonning vm: %s", err)
	}

	vm_mor := info.Result.(types.ManagedObjectReference)
	d.SetId(vm_mor.Value)
	vm := object.NewVirtualMachine(client, vm_mor)
	// workaround for https://github.com/vmware/govmomi/issues/218
	if ip_address == "" && d.Get("power_on").(bool) {
		ip, err := vm.WaitForIP(ctx)
		if err != nil {
			log.Printf("[ERROR] Cannot read ip address: %s", err)
		} else {
			d.Set("ip_address", ip)
			d.SetConnInfo(map[string]string{
				"type": "ssh",
				"host": ip,
			})
		}
	}

	return nil
}

func resourceVirtualMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	ctx := context.TODO()
	vm_mor := types.ManagedObjectReference{Type: "VirtualMachine", Value: d.Id()}
	vm := object.NewVirtualMachine(client, vm_mor)

	var vm_mo mo.VirtualMachine
	err := vm.Properties(ctx, vm.Reference(), []string{"summary"}, &vm_mo)
	if err != nil {
		log.Printf("[INFO] Cannot read VM properties: %s", err)
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
		ip, err := vm.WaitForIP(ctx)
		if err != nil {
			log.Printf("[ERROR] Cannot read ip address: %s", err)
		} else {
			d.Set("ip_address", ip)
		}
	}

	return nil
}

func resourceVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	ctx := context.TODO()

	vm_mor := types.ManagedObjectReference{Type: "VirtualMachine", Value: d.Id()}
	vm := object.NewVirtualMachine(client, vm_mor)

	task, err := vm.PowerOff(ctx)
	if err != nil {
		return fmt.Errorf("Error powering vm off: %s", err)
	}
	task.WaitForResult(ctx, nil)

	task, err = vm.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting vm: %s", err)
	}
	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		return fmt.Errorf("Error deleting vm: %s", err)
	}

	return nil
}
