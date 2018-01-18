package vsphere

import (
	"context"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"log"
	"strings"
)

func resourceVirtualMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineCreate,
		Read:   resourceVirtualMachineRead,
		Delete: resourceVirtualMachineDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"folder": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"resource_pool": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"datastore": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"linked_clone": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
			"cpus": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     resourceVirtualMachineDisk(),
			},
			"domain": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"subnet_mask": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"configuration_parameters": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"power_on": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
				ForceNew: true,
			},
		},
	}
}

func resourceVirtualMachineCreate(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(providerMeta)
	client := providerMeta.client
	ctx := providerMeta.context
	finder := find.NewFinder(client, false)

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
			return errors.New("Image is a template, 'host' is a required")
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

	var devices object.VirtualDeviceList

	for _, diskValue := range d.Get("disks").([]interface{}) {
		if disk, ok := diskValue.(map[string]interface{}); ok {
			diskDatastoreName := disk["datastore"].(string)

			datastore, err := finder.Datastore(ctx, fmt.Sprintf("/%v/datastore/%v", dc_name, diskDatastoreName))
			if err != nil {
				return fmt.Errorf("Failed to find datastore \"%s\": %v", diskDatastoreName, err)
			}

			datastoreRef := datastore.Reference()

			controllerDevice, err := devices.CreateSCSIController("")
			if err != nil {
				return err
			}

			devices = append(devices, controllerDevice)

			diskDevice := &types.VirtualDisk{
				VirtualDevice: types.VirtualDevice{
					Key: devices.NewKey(),
					Backing: &types.VirtualDiskFlatVer2BackingInfo{
						DiskMode: string(types.VirtualDiskModeIndependent_persistent),
						VirtualDeviceFileBackingInfo: types.VirtualDeviceFileBackingInfo{
							FileName:  datastore.Path(ensureVmdkSuffix(disk["path"].(string))),
							Datastore: &datastoreRef,
						},
					},
				},
			}

			devices = append(devices, diskDevice)

			devices.AssignController(diskDevice, controllerDevice.(types.BaseVirtualController))

			confSpec.DeviceChange = append(confSpec.DeviceChange, &types.VirtualDeviceConfigSpec{
				Operation: types.VirtualDeviceConfigSpecOperationAdd,
				Device:    controllerDevice,
			}, &types.VirtualDeviceConfigSpec{
				Operation: types.VirtualDeviceConfigSpecOperationAdd,
				Device:    diskDevice,
			})
		}
	}

	params := d.Get("configuration_parameters").(map[string]interface{})
	var ov []types.BaseOptionValue
	if len(params) > 0 {
		for k, v := range params {
			o := types.OptionValue{
				Key:   k,
				Value: v,
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
			return errors.New("`linked_clone=true`, but image VM has no snapshots")
		}
		cloneSpec.Snapshot = image_mo.Snapshot.CurrentSnapshot
	}

	domain := d.Get("domain").(string)
	ip_address := d.Get("ip_address").(string)
	if domain != "" {
		if image_mo.Guest.ToolsVersionStatus2 == "guestToolsNotInstalled" {
			return errors.New("VMware tools are not installed in base VM")
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
				return errors.New("'subnet_mask' must be set, if static 'ip_address' is specified")
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
		return errors.New("'domain' must be set, if static 'ip_address' is specified")
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
	if d.Get("power_on").(bool) {
		if ip_address == "" {
			ip, err := updateIPAddress(vm, ctx)
			if err != nil {
				return fmt.Errorf("[ERROR] Cannot read ip addresses: %s", err)
			}
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
	providerMeta := meta.(providerMeta)
	client := providerMeta.client
	ctx := providerMeta.context
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
		ip, err := updateIPAddress(vm, ctx)
		if err != nil {
			return fmt.Errorf("[ERROR] Cannot read ip addresses: %s", err)
		}
		d.Set("ip_address", ip)
	}

	return nil
}

func resourceVirtualMachineDelete(d *schema.ResourceData, meta interface{}) error {
	providerMeta := meta.(providerMeta)
	client := providerMeta.client
	ctx := providerMeta.context

	vm_mor := types.ManagedObjectReference{Type: "VirtualMachine", Value: d.Id()}
	vm := object.NewVirtualMachine(client, vm_mor)

	task, err := vm.PowerOff(ctx)
	if err != nil {
		return fmt.Errorf("Error powering vm off: %s", err)
	}
	err = task.Wait(ctx)
	if err != nil {
		if !strings.Contains(err.Error(), "in the current state (Powered off)") {
			return fmt.Errorf("Error powering vm off: %s", err)
		}
	}

	devices, err := vm.Device(ctx)
	if err != nil {
		return err
	}

	var devicesToDestroy object.VirtualDeviceList
	for _, device := range devices {
		if disk, ok := device.(*types.VirtualDisk); ok {
			if diskBacking, ok := disk.Backing.(*types.VirtualDiskFlatVer2BackingInfo); ok {
				// detach all independent persistent disks,
				// otherwise all dynamically attached disks are destroyed along with VM losing its data
				if diskBacking.DiskMode == string(types.VirtualDiskModeIndependent_persistent) {
					devicesToDestroy = append(devicesToDestroy, &types.VirtualDevice{
						Key: device.GetVirtualDevice().Key,
					})
				}
			}
		}
	}

	if err := vm.RemoveDevice(ctx, true, devicesToDestroy...); err != nil {
		return err
	}

	task, err = vm.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting vm: %s", err)
	}
	err = task.Wait(ctx)
	if err != nil {
		return fmt.Errorf("Error deleting vm: %s", err)
	}

	return nil
}

func updateIPAddress(vm *object.VirtualMachine, ctx context.Context) (string, error) {
	ipMap, err := vm.WaitForNetIP(ctx, true)
	if err != nil {
		return "", err
	}
	if len(ipMap) > 0 {
		for _, ip := range ipMap {
			return ip[0], nil
		}
	}

	return "", nil
}
