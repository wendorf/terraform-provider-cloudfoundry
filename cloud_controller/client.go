package cloud_controller

import (
	"fmt"
	"golang.org/x/oauth2"
	"net"
	"time"
	"net/http"
	"crypto/tls"
	"context"
	"encoding/json"
	"io/ioutil"
	"bytes"
)

type Client struct {
	Organizations OrganizationsClient
	Spaces        SpacesClient
	Services      ServicesClient
	ServicePlans  ServicePlansClient
	Domains       DomainsClient
	httpClient    *http.Client
	apiEndpoint   string
}

type info struct {
	AuthorizationEndpoint string `json:"authorization_endpoint"`
	TokenEndpoint         string `json:"token_endpoint"`
}

func NewClient(apiEndpoint, username, password string) (*Client, error) {
	info, err := newInfo(apiEndpoint)
	if err != nil {
		return nil, err
	}

	oauthClient, err := newOauthClient(info.AuthorizationEndpoint, info.TokenEndpoint, username, password)
	if err != nil {
		return nil, err
	}

	client := &Client{
		httpClient: oauthClient,
		apiEndpoint: apiEndpoint,
	}
	client.Organizations = OrganizationsClient{client: client, }
	client.Spaces = SpacesClient{client: client, }
	client.Services = ServicesClient{client: client, }
	client.ServicePlans = ServicePlansClient{client: client, }
	client.Domains = DomainsClient{client: client, }
	return client, nil
}

func (c *Client) Get(path string) (*http.Response, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", c.apiEndpoint, path), nil)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) Delete(path string) (*http.Response, error) {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s%s", c.apiEndpoint, path), nil)
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) Post(path string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", c.apiEndpoint, path), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func (c *Client) Put(path string, body interface{}) (*http.Response, error) {
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s", c.apiEndpoint, path), bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}

	return c.httpClient.Do(req)
}

func newInfo(apiEndpoint string) (info, error) {
	var (
		info info
		response *http.Response
		err error
		body []byte
	)

	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}}

	if response, err = client.Get(fmt.Sprintf("%s/v2/info", apiEndpoint)); err != nil {
		return info, err
	}

	if body, err = ioutil.ReadAll(response.Body); err != nil {
		return info, err
	}

	if err = json.Unmarshal(body, &info); err != nil {
		return info, err
	}

	return info, nil
}

func newOauthClient(authorizationEndpoint, tokenEndpoint, username, password string) (*http.Client, error) {
	conf := &oauth2.Config{
		ClientID:     "cf",
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			AuthURL: fmt.Sprintf("%s/oauth/authorize", authorizationEndpoint),
			TokenURL: fmt.Sprintf("%s/oauth/token", tokenEndpoint),
		},
	}

	httpClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			Dial: (&net.Dialer{
				Timeout:   5 * time.Second,
				KeepAlive: 30 * time.Second,
			}).Dial,
			ResponseHeaderTimeout: 15 * time.Minute,
		},
		Timeout: 30 * time.Minute,
	}

	insecureContext := context.Background()
	insecureContext = context.WithValue(insecureContext, oauth2.HTTPClient, httpClient)

	token, err := conf.PasswordCredentialsToken(insecureContext, username, password)
	if err != nil {
		return nil, err
	}

	return conf.Client(insecureContext, token), nil
}