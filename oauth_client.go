package main

import (
	"fmt"
	"golang.org/x/oauth2"
	"net"
	"time"
	"net/http"
	"crypto/tls"
	"context"
)

func NewOAuthHTTPClient(host, username, password string) (*http.Client, error) {
	conf := &oauth2.Config{
		ClientID:     "cf",
		ClientSecret: "",
		Endpoint: oauth2.Endpoint{
			TokenURL: fmt.Sprintf("%s/oauth/token", host),
		},
	}

	httpclient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
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
	insecureContext = context.WithValue(insecureContext, oauth2.HTTPClient, httpclient)

	token, err := conf.PasswordCredentialsToken(insecureContext, username, password)
	if err != nil {
		return nil, err
	}

	return conf.Client(insecureContext, token), nil
}