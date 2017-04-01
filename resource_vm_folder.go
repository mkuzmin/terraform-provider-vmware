package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"golang.org/x/net/context"
	"fmt"
)

func resourceVmFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceVmFolderCreate,
		Read:   resourceVmFolderRead,
		Delete: resourceVmFolderDelete,

		Schema: map[string]*schema.Schema{
			"parent": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			// create parent folders
		},
	}
}

func resourceVmFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	finder := find.NewFinder(client, false)
	ctx := context.TODO()

	parent_name := d.Get("parent").(string)
	name := d.Get("name").(string)

	parent_folder, err := finder.Folder(ctx, parent_name)
	if err != nil {
		return fmt.Errorf("Cannot find parent folder: %s", err)
	}

	folder, err := parent_folder.CreateFolder(ctx, name)
	if err != nil {
		return fmt.Errorf("Cannot create folder: %s", err)
	}

	d.SetId(folder.InventoryPath)

	return nil
}

func resourceVmFolderRead(d *schema.ResourceData, meta interface{}) error {

	return nil
}

func resourceVmFolderDelete(d *schema.ResourceData, meta interface{}) error {

	// check VMs inside

	return nil
}
