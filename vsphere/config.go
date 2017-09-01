package vsphere

import (
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25"
	"golang.org/x/net/context"
	"log"
	"net/url"
)

type Config struct {
	vCenter  string
	User     string
	Password string
	Insecure bool
}

func (c *Config) Client(ctx context.Context) (*vim25.Client, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", c.User, c.Password, c.vCenter))
	if err != nil {
		return nil, fmt.Errorf("Incorrect vCenter server address: %s", err)
	}
	client, err := govmomi.NewClient(ctx, u, c.Insecure)
	if err != nil {
		return nil, fmt.Errorf("Error setting up client: %s", err)
	}
	log.Println("[INFO] vSphere Client configured")
	return client.Client, nil
}
