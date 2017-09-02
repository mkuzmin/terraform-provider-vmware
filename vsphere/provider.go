package vsphere

import (
	"context"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/vmware/govmomi/vim25"
)

type providerMeta struct {
	context context.Context
	client  *vim25.Client
}

func Provider() terraform.ResourceProvider {
	provider := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"vcenter_server": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_SERVER", nil),
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
				DefaultFunc: schema.EnvDefaultFunc("VSPHERE_INSECURE", false),
				Description: "Do not check vCenter SSL certificate",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"vmware_virtual_disk": resourceVirtualDisk(),
			"vmware_virtual_machine": resourceVirtualMachine(),
			"vmware_vm_folder":       resourceVmFolder(),
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

		return providerMeta{
			client:  client,
			context: ctx,
		}, err

	}
}
