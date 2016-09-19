package main

import (
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller"
)

type Config struct {
	ApiEndpoint      string
	Username         string
	Password         string
	OrganizationGUID string
	Client           *cloud_controller.Client
}

func (c *Config) load() error {
	client, err := cloud_controller.NewClient(c.ApiEndpoint, c.Username, c.Password)
	if err != nil {
		return err
	}
	c.Client = client
	return nil
}