package main

import (
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVirtualMachineDisk() *schema.Resource {
	return &schema.Resource{
		Create: resourceVirtualMachineDiskCreate,
		Read:   resourceVirtualMachineDiskRead,
		Delete: resourceVirtualMachineDiskDelete,

		Schema: map[string]*schema.Schema{
			"datastore": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceVirtualMachineDiskCreate(resourceData *schema.ResourceData, _ interface{}) error {
	id, err := uuid.GenerateUUID()
	if err != nil {
		return err
	}

	resourceData.SetId(id)

	return nil
}

func resourceVirtualMachineDiskRead(_ *schema.ResourceData, _ interface{}) error {
	return nil
}

func resourceVirtualMachineDiskDelete(resourceData *schema.ResourceData, _ interface{}) error {
	resourceData.SetId("")

	return nil
}
