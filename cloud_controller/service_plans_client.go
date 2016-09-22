package cloud_controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"encoding/json"
)

type ServicePlansClient struct {
	client *Client
}

func (c *ServicePlansClient) GetGUID(service, servicePlan string) (string, error) {
	services, err := c.client.Services.List(fmt.Sprintf("label:%s", service))
	if err != nil {
		return "", err
	}

	servicePlanCollection, err := c.List(fmt.Sprintf("service_guid:%s", services.Services[0].Metadata.GUID))
	if err != nil {
		return "", err
	}

	for _, servicePlanWrapper := range servicePlanCollection.ServicePlans {
		if servicePlanWrapper.ServicePlan.Name == servicePlan {
			return servicePlanWrapper.Metadata.GUID, nil
		}
	}

	return "", errors.New(fmt.Sprintf("Could not find service plan %s for service %s", servicePlan, service))
}

func (c *ServicePlansClient) List(query string) (models.ServicePlanCollection, error) {
	var servicePlanCollection models.ServicePlanCollection

	resp, err := c.client.Get(fmt.Sprintf("/v2/service_plans?q=%s", query))
	if err != nil {
		return servicePlanCollection, err
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return servicePlanCollection, errors.New(fmt.Sprintf("Could not fetch service plans: %s\n%s", query, string(responseBody)))
	}

	json.Unmarshal(responseBody, &servicePlanCollection)

	return servicePlanCollection, nil
}