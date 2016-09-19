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

func (c *Config) load(organizationName string) error {
	client, err := cloud_controller.NewClient(c.ApiEndpoint, c.Username, c.Password)
	if err != nil {
		return err
	}
	c.Client = client

	c.OrganizationGUID, err = c.Client.Organizations.GetGUID(organizationName)
	if err != nil {
		return err
	}

	return nil
}