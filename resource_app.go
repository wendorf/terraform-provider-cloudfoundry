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

func resourceApp() *schema.Resource {
	return &schema.Resource{
		Create: resourceAppCreate,
		Read:   resourceAppRead,
		Update: resourceAppUpdate,
		Delete: resourceAppDelete,

		Schema: map[string]*schema.Schema{
			"name": {
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

func resourceAppCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	spaceGUID, err := config.Client.Spaces.GetGUID(d.Get("space").(string))
	if err != nil {
		return err
	}

	app := models.App{
		Name: d.Get("name").(string),
		Relationships: models.AppRelationships{
			Space: models.GUID{
				GUID: spaceGUID,
			},
		},
	}

	resp, err := config.Client.Post("/v3/apps", app)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not create app %s: %s", d.Get("name").(string), string(responseBody)))
	}

	json.Unmarshal(responseBody, &app)

	d.SetId(app.GUID)
	return nil
}

func resourceAppRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Get(fmt.Sprintf("/v3/apps/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		d.SetId("")
	}

	return nil
}

func resourceAppUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceAppDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}