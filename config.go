package main

import (
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller"
	"fmt"
	"net/http"
	"errors"
	"io/ioutil"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"encoding/json"
)

type Config struct {
	ApiEndpoint      string
	Username         string
	Password         string
	OrganizationGUID string
	Client           *cloud_controller.Client
}

func (c *Config) load(organizationName string) error {
	client, err := cloud_controller.NewClient(c.ApiEndpoint, c.Username, c.Password)
	if err != nil {
		return err
	}
	c.Client = client

	resp, err := c.Client.Get(fmt.Sprintf("/v2/organizations?q=name:%s", organizationName))
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("Could not fetch organization %s", organizationName))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	v2Collection := models.V2Collection{}
	json.Unmarshal(responseBody, &v2Collection)

	c.OrganizationGUID = v2Collection.V2Objects[0].Metadata.GUID

	return nil
}