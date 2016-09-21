package cloud_controller

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"errors"
	"github.com/wendorf/terraform-provider-cloudfoundry/cloud_controller/models"
	"encoding/json"
)

type SpacesClient struct {
	client  *Client
}

func (c *SpacesClient) GetGUID(name string) (string, error) {
	spaceCollection, err := c.List(fmt.Sprintf("name:%s", name))
	if err != nil {
		return "", err
	}

	return spaceCollection.Spaces[0].Metadata.GUID, nil
}

func (c *SpacesClient) List(query string) (models.SpaceCollection, error) {
	var spaceCollection models.SpaceCollection

	resp, err := c.client.Get(fmt.Sprintf("/v2/spaces?q=%s", query))
	if err != nil {
		return spaceCollection, err
	}

	if resp.StatusCode != http.StatusOK {
		return spaceCollection, errors.New(fmt.Sprintf("Could not fetch spaces: %s", query))
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(responseBody, &spaceCollection)

	return spaceCollection, nil
}