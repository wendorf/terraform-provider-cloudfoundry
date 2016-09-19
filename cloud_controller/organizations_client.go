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
	resp, err := c.client.Get(fmt.Sprintf("/v2/organizations?q=name:%s", name))
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", errors.New(fmt.Sprintf("Could not fetch organization %s", name))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	v2Collection := models.V2Collection{}
	json.Unmarshal(responseBody, &v2Collection)

	return v2Collection.V2Objects[0].Metadata.GUID, nil
}