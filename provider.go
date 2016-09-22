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
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"skip_ssl_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default: false,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"cloudfoundry_space": resourceSpace(),
			"cloudfoundry_service_instance": resourceServiceInstance(),
			"cloudfoundry_route": resourceRoute(),
			"cloudfoundry_app": resourceApp(),
		},
		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		ApiEndpoint: d.Get("api_endpoint").(string),
		Username:     d.Get("username").(string),
		Password:      d.Get("password").(string),
		SkipSSLValidation:      d.Get("skip_ssl_validation").(bool),
	}

	if err := config.load(d.Get("organization").(string)); err != nil {
		return nil, err
	}

	return &config, nil
}