package main

import (
	"context"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

type providerMeta struct {
	context context.Context
	client  *vim25.Client
	folders *object.DatacenterFolders
}

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
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
			"datacenter": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmware_virtual_machine": resourceVirtualMachine(),
		},
	}
	provider.ConfigureFunc = providerConfigure(provider)

	return provider
}

func providerConfigure(p *schema.Provider) schema.ConfigureFunc {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := Config{
			vCenter:  d.Get("vcenter_server").(string),
			User:     d.Get("user").(string),
			Password: d.Get("password").(string),
			Insecure: d.Get("insecure_connection").(bool),
		}

		ctx := p.StopContext()
		client, err := config.Client(ctx)
		finder := find.NewFinder(client, false)

		datacenter := d.Get("datacenter").(string)
		var dc *object.Datacenter
		if datacenter == "" {
			dc, err = finder.DefaultDatacenter(ctx)
			if err != nil {
				return nil, err
			}
		} else {
			dc, err = finder.Datacenter(ctx, datacenter)
			if err != nil {
				return nil, err
			}
		}
		folders, err := dc.Folders(ctx)
		if err != nil {
			return nil, err
		}

		return providerMeta{
			client:  client,
			context: ctx,
			folders: folders,
		}, nil

	}
}
