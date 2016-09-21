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

	services, err := config.Client.Services.List(fmt.Sprintf("label:%s", d.Get("service").(string)))
	if err != nil {
		return err
	}

	servicePlanCollection, err := config.Client.ServicePlans.List(fmt.Sprintf("service_guid:%s", services.Services[0].Metadata.GUID))
	if err != nil {
		return err
	}

	var servicePlanGUID string
	for _, servicePlanWrapper := range servicePlanCollection.ServicePlans {
		if servicePlanWrapper.ServicePlan.Name == d.Get("service_plan").(string) {
			servicePlanGUID = servicePlanWrapper.Metadata.GUID
		}
	}

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

	if resp.StatusCode != http.StatusCreated {
		return errors.New(fmt.Sprintf("Could not create service instance %s: %s", d.Get("name").(string), string(responseBody)))
	}

	serviceInstanceWrapper := models.ServiceInstanceWrapper{}
	json.Unmarshal(responseBody, &serviceInstanceWrapper)

	d.SetId(serviceInstanceWrapper.Metadata.GUID)

	return nil
}

func resourceServiceInstanceRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceInstanceUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceServiceInstanceDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
