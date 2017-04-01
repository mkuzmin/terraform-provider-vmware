package main

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	"golang.org/x/net/context"
	"strings"
)

func resourceVirtualDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualDiskCreate,
		Read:   resourceVirtualDiskRead,
		Delete: resourceVirtualDiskDelete,

		Schema: map[string]*schema.Schema{
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"datastore": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sizeMb": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"thick": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
				ForceNew: true,
			},
		},
	}
}

func resourceVirtualDiskCreate(resourceData *schema.ResourceData, meta interface{}) error {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	client := meta.(*vim25.Client)
	finder := find.NewFinder(client, false)
	ctx := context.TODO()

	datacenter, err := findDatacenter(ctx, finder, resourceData)
	if err != nil {
		return err
	}

	finder.SetDatacenter(datacenter)

	datastore, err := findDatastore(ctx, finder, resourceData)
	if err != nil {
		return err
	}

	diskPath, err := getDiskPath(resourceData, datastore)
	if err != nil {
		return err
	}

	diskSizeMb := resourceData.Get("sizeMb").(int)
	if diskSizeMb == 0 {
		return errors.New("Virtual disk size is not specified")
	}

	diskType := string(types.VirtualDiskTypeThin)
	if resourceData.Get("thick").(bool) {
		diskType = string(types.VirtualDiskTypeThick)
	}

	diskManager := object.NewVirtualDiskManager(client)

	diskTask, err := diskManager.CreateVirtualDisk(ctx, diskPath, datacenter, &types.FileBackedVirtualDiskSpec{
		VirtualDiskSpec: types.VirtualDiskSpec{
			DiskType:    diskType,
			AdapterType: string(types.VirtualDiskAdapterTypeLsiLogic),
		},
		CapacityKb: int64(diskSizeMb) * 1024,
	})

	if err != nil {
		return fmt.Errorf("Failed to create virtual disk: %v", err)
	}

	if err := diskTask.Wait(ctx); err != nil {
		return fmt.Errorf("Failed to create virtual disk: %v", err)
	}

	resourceData.SetId(id)

	return nil
}

func resourceVirtualDiskRead(_ *schema.ResourceData, _ interface{}) error {
    return nil // todo
}

func resourceVirtualDiskDelete(resourceData *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	finder := find.NewFinder(client, false)
	ctx := context.TODO()

	datacenter, err := findDatacenter(ctx, finder, resourceData)
	if err != nil {
		return err
	}

	finder.SetDatacenter(datacenter)

	datastore, err := findDatastore(ctx, finder, resourceData)
	if err != nil {
		return err
	}

	diskPath, err := getDiskPath(resourceData, datastore)
	if err != nil {
		return err
	}

	diskManager := object.NewVirtualDiskManager(client)

	diskTask, err := diskManager.DeleteVirtualDisk(ctx, diskPath, datacenter)
	if err != nil {
		return fmt.Errorf("Failed to destroy virtual disk: %v", err)
	}

	if err := diskTask.Wait(ctx); err != nil {
		return fmt.Errorf("Failed to destroy virtual disk: %v", err)
	}

	resourceData.SetId("")

	return nil
}

func findDatacenter(ctx context.Context, finder *find.Finder, resourceData *schema.ResourceData) (*object.Datacenter, error) {
	datacenterName := resourceData.Get("datacenter").(string)
	if datacenterName == "" {
		datacenter, err := finder.DefaultDatacenter(ctx)
		if err != nil {
			return nil, fmt.Errorf("Failed to read default datacenter: %v", err)
		}

		var moDatacenter *mo.Datacenter
		if err = datacenter.Properties(ctx, datacenter.Reference(), []string{"name"}, moDatacenter); err != nil {
			return nil, fmt.Errorf("Failed to read default datacenter name: %v", err)
		}

		resourceData.Set("datacenter", moDatacenter.Name)

		return datacenter, nil
	}

	datacenter, err := finder.Datacenter(ctx, datacenterName)
	if err != nil {
		return nil, fmt.Errorf("Failed to find datacenter \"%s\": %v", datacenterName, err)
	}

	return datacenter, nil
}

func findDatastore(ctx context.Context, finder *find.Finder, resourceData *schema.ResourceData) (*object.Datastore, error) {
	datastoreName := resourceData.Get("datastore").(string)
	if datastoreName == "" {
		datastore, err := finder.DefaultDatastore(ctx)
		if err != nil {
			return nil, fmt.Errorf("Failed to read default datastore: %v", err)
		}

		var moDatastore *mo.Datastore
		if err = datastore.Properties(ctx, datastore.Reference(), []string{"name"}, moDatastore); err != nil {
			return nil, fmt.Errorf("Failed to read default datastore name: %v", err)
		}

		resourceData.Set("datastore", moDatastore.Name)

		return datastore, nil
	}

	datastore, err := finder.Datastore(ctx, datastoreName)
	if err != nil {
		return nil, fmt.Errorf("Failed to find datastore \"%s\": %v", datastore, err)
	}

	return datastore, nil
}

func getDiskPath(resourceData *schema.ResourceData, datastore *object.Datastore) (string, error) {
	diskPath := resourceData.Get("path").(string)
	if diskPath == "" {
		return "", errors.New("Virtual disk path is not specified")
	}

	if !strings.HasSuffix(diskPath, ".vmdk") {
		diskPath = fmt.Sprintf("%s.vmdk", diskPath)
	}

	return datastore.Path(diskPath), nil
}
