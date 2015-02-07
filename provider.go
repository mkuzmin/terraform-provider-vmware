package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vcenter_server": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "vCenter server address",
			},
			"user": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USER", nil),
				Description: "User account",
			},
			"password": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
				Description: "Password",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vsphere_virtual_machine": resourceVirtualMachine(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		vCenter: d.Get("vcenter_server").(string),
		User: d.Get("user").(string),
		Password: d.Get("password").(string),
	}

	return config.Client()
}