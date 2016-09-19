package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"api_endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"username": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_guid": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudfoundry_space": resourceSpace(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ApiEndpoint: d.Get("api_endpoint").(string),
		Username:     d.Get("username").(string),
		Password:      d.Get("password").(string),
		OrganizationGUID:      d.Get("organization_guid").(string),
	}

	if err := config.load(); err != nil {
		return nil, err
	}

	return &config, nil
}