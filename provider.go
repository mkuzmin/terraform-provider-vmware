package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vcenter_server": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "vCenter server address",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_USER", nil),
				Description: "User account",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_PASSWORD", nil),
				Description: "Password",
				Sensitive:   true,
			},
			"insecure_connection": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Do not check vCenter SSL certificate",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmware_virtual_machine": resourceVirtualMachine(),
			"vmware_vm_folder":       resourceVmFolder(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		vCenter:  d.Get("vcenter_server").(string),
		User:     d.Get("user").(string),
		Password: d.Get("password").(string),
		Insecure: d.Get("insecure_connection").(bool),
	}

	return config.Client()
}
