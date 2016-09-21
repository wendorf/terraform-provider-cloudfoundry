package cloud_controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"encoding/json"
)

type ServicesClient struct {
	client  *Client
}

func (c *ServicesClient) List(query string) (models.ServiceCollection, error) {
	var serviceCollection models.ServiceCollection

	resp, err := c.client.Get(fmt.Sprintf("/v2/services?q=%s", query))
	if err != nil {
		return serviceCollection, err
	}

	if resp.StatusCode != http.StatusOK {
		return serviceCollection, errors.New(fmt.Sprintf("Could not fetch services: %s", query))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseBody, &serviceCollection)

	return serviceCollection, nil
}