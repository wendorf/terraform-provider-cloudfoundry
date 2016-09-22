package cloud_controller

import (
"fmt"
"net/http"
"io/ioutil"
"errors"
"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
"encoding/json"
)

type DomainsClient struct {
	client  *Client
}

func (c *DomainsClient) GetGUID(name string) (string, error) {
	domainCollection, err := c.List(fmt.Sprintf("name:%s", name))
	if err != nil {
		return "", err
	}

	return domainCollection.Domains[0].Metadata.GUID, nil
}

func (c *DomainsClient) List(query string) (models.DomainCollection, error) {
	var domainCollection models.DomainCollection

	resp, err := c.client.Get(fmt.Sprintf("/v2/domains?q=%s", query))
	if err != nil {
		return domainCollection, err
	}

	if resp.StatusCode != http.StatusOK {
		return domainCollection, errors.New(fmt.Sprintf("Could not fetch domains: %s", query))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseBody, &domainCollection)

	return domainCollection, nil
}