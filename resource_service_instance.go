package main

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"fmt"
	"io/ioutil"
	"net/http"
	"errors"
	"encoding/json"
)

func resourceServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceServiceInstanceCreate,
		Read:   resourceServiceInstanceRead,
		Update: resourceServiceInstanceUpdate,
		Delete: resourceServiceInstanceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"space": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_plan": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceServiceInstanceCreate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	spaceGUID, err := config.Client.Spaces.GetGUID(d.Get("space").(string))
	if err != nil {
		return err
	}

	servicePlanGUID, err := config.Client.ServicePlans.GetGUID(d.Get("service").(string), d.Get("service_plan").(string))

	serviceInstance := models.ServiceInstance{
		Name: d.Get("name").(string),
		SpaceGUID: spaceGUID,
		ServicePlanGUID: servicePlanGUID,
	}

	resp, err := config.Client.Post("/v2/service_instances", serviceInstance)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("Could not create service instance %s: %s", d.Get("name").(string), string(responseBody)))
	}

	serviceInstanceWrapper := models.ServiceInstanceWrapper{}
	json.Unmarshal(responseBody, &serviceInstanceWrapper)

	d.SetId(serviceInstanceWrapper.Metadata.GUID)

	return nil
}

func resourceServiceInstanceRead(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Get(fmt.Sprintf("/v2/service_instances/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		d.SetId("")
	}

	return nil
}

func resourceServiceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)

	servicePlanGUID, err := config.Client.ServicePlans.GetGUID(d.Get("service").(string), d.Get("service_plan").(string))

	serviceInstance := models.ServiceInstance{
		Name: d.Get("name").(string),
		ServicePlanGUID: servicePlanGUID,
	}

	resp, err := config.Client.Put(fmt.Sprintf("/v2/service_instances/%s", d.Id()), serviceInstance)
	if err != nil {
		return err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not create service instance %s: %s", d.Get("name").(string), string(responseBody)))
	}

	serviceInstanceWrapper := models.ServiceInstanceWrapper{}
	json.Unmarshal(responseBody, &serviceInstanceWrapper)

	d.SetId(serviceInstanceWrapper.Metadata.GUID)

	return nil
}

func resourceServiceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(*Config)
	resp, err := config.Client.Delete(fmt.Sprintf("/v2/service_instances/%s", d.Id()))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New(fmt.Sprintf("Could not delete service instance %s", d.Get("name").(string)))
	}

	return nil
}
