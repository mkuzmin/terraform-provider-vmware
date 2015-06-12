package main

import (
	"fmt"
	"log"
	"net/url"
	"github.com/vmware/govmomi"
	"golang.org/x/net/context"
)

type Config struct {
	vCenter string
	User string
	Password string
    Insecure bool
}

func (c *Config) Client() (*govmomi.Client, error) {
	u, err := url.Parse(fmt.Sprintf ("https://%s:%s@%s/sdk", c.User, c.Password, c.vCenter))
	if err != nil {
		return nil, fmt.Errorf("Incorrect vCenter server address: %s", err)
	}
	client, err := govmomi.NewClient(context.TODO(), u, c.Insecure)
	if err != nil {
		return nil, fmt.Errorf("Error setting up client: %s", err)
	}
	log.Printf("[INFO] vSphere Client configured")
	return client, nil
}
