package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/vim25"
	"golang.org/x/net/context"
	"fmt"
	"github.com/vmware/govmomi/vim25/types"
	"github.com/vmware/govmomi/object"
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

	d.SetId(folder.Reference().Value)
	return nil
}

func resourceVmFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	finder := find.NewFinder(client, false)
	ctx := context.TODO()

	mor := types.ManagedObjectReference{Type: "Folder", Value: d.Id()}
	obj, err := finder.ObjectReference(ctx, mor)
	if err != nil {
		d.SetId("")
		return nil
	}

	folder := obj.(*object.Folder)
	name, err := folder.ObjectName(ctx)
	if err != nil {
		return fmt.Errorf("Cannot read folder: %s", err)
	}

	d.Set("name", name)
	return nil
}

func resourceVmFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*vim25.Client)
	finder := find.NewFinder(client, false)
	ctx := context.TODO()

	mor := types.ManagedObjectReference{Type: "Folder", Value: d.Id()}
	obj, err := finder.ObjectReference(ctx, mor)
	if err != nil {
		d.SetId("")
		return nil
	}
	folder := obj.(*object.Folder)

	if children, _ := folder.Children(ctx); len(children) > 0 {
		return fmt.Errorf("Folder is not empty")
	}

	task, err := folder.Destroy(ctx)
	if err != nil {
		return fmt.Errorf("Cannot delete folder: %s", err)
	}
	_, err = task.WaitForResult(ctx, nil)
	if err != nil {
		return fmt.Errorf("Cannot delete folder: %s", err)
	}

	d.SetId("")
	return nil
}
