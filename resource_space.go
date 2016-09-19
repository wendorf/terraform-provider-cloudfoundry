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

func resourceSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpaceCreate,
		Read:   resourceSpaceRead,
		Update: resourceSpaceUpdate,
		Delete: resourceSpaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSpaceCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	space := models.Space{
		Name: d.Get("name").(string),
		OrganizationGUID: config.OrganizationGUID,
	}
	resp, err := config.Client.Post("/v2/spaces", space)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not create space %s: %s", d.Get("name").(string), string(responseBody)))
	}

	v2Object := models.V2Object{}
	json.Unmarshal(responseBody, &v2Object)

	d.SetId(v2Object.Metadata.GUID)
	return nil
}

func resourceSpaceRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Get(fmt.Sprintf("/v2/spaces/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		d.SetId("")
	}

	return nil
}

func resourceSpaceUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	space := models.Space{
		Name: d.Get("name").(string),
		OrganizationGUID: config.OrganizationGUID,
	}
	resp, err := config.Client.Put(fmt.Sprintf("/v2/spaces/%s", d.Id()), space)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not update space %s", d.Get("name").(string)))
	}

	return nil
}

func resourceSpaceDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Delete(fmt.Sprintf("/v2/spaces/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Could not delete space %s", d.Get("name").(string)))
	}

	return nil
}
