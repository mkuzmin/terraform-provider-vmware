package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceVm() *schema.Resource {
	return &schema.Resource{
		Create: resourceVmCreate,
		Read:   resourceVmRead,
		Update: resourceVmUpdate,
		Delete: resourceVmDelete,

		Schema: map[string]*schema.Schema{
		},
	}
}

func resourceVmCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVmRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVmUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceVmDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
