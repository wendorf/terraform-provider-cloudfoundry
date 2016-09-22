package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"net/http"
	"errors"
)

func resourceRoute() *schema.Resource {
	return &schema.Resource{
		Create: resourceRouteCreate,
		Read:   resourceRouteRead,
		Update: resourceRouteUpdate,
		Delete: resourceRouteDelete,

		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Required: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},
			"space": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceRouteCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	spaceGuid, err := config.Client.Spaces.GetGUID(d.Get("space").(string))
	if err != nil {
		return err
	}

	domainGuid, err := config.Client.Domains.GetGUID(d.Get("domain").(string))
	if err != nil {
		return err
	}

	route := models.Route{
		Host: d.Get("host").(string),
		SpaceGUID: spaceGuid,
		DomainGUID: domainGuid,
	}

	resp, err := config.Client.Post("/v2/routes", route)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not create route %s: %s", d.Get("host").(string), string(responseBody)))
	}

	var domainWrapper models.DomainWrapper
	json.Unmarshal(responseBody, &domainWrapper)

	d.SetId(domainWrapper.Metadata.GUID)
	return nil
}

func resourceRouteRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Get(fmt.Sprintf("/v2/routes/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		d.SetId("")
	}

	return nil
}

func resourceRouteUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceRouteDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Delete(fmt.Sprintf("/v2/routes/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Could not delete route %s", d.Get("host").(string)))
	}

	return nil
}