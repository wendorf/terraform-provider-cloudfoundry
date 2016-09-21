package cloud_controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"encoding/json"
)

type OrganizationsClient struct {
	client  *Client
}

func (c *OrganizationsClient) GetGUID(name string) (string, error) {
	organizationCollection, err := c.List(fmt.Sprintf("name:%s", name))
	if err != nil {
		return "", err
	}

	return organizationCollection.Organizations[0].Metadata.GUID, nil
}

func (c *OrganizationsClient) List(query string) (models.OrganizationCollection, error) {
	var organizationCollection models.OrganizationCollection

	resp, err := c.client.Get(fmt.Sprintf("/v2/organizations?q=%s", query))
	if err != nil {
		return organizationCollection, err
	}

	if resp.StatusCode != http.StatusOK {
		return organizationCollection, errors.New(fmt.Sprintf("Could not fetch organizations: %s", query))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseBody, &organizationCollection)

	return organizationCollection, nil
}