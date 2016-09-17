package main

import (
    "github.com/hashicorp/terraform/plugin"
    "github.com/hashicorp/terraform/terraform"
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {
    _, err := getSampleResponse("username", "password")
    if err != nil {
        fmt.Printf("Could not get sample request: %s", err)
    }

    plugin.Serve(&plugin.ServeOpts{
        ProviderFunc: func() terraform.ResourceProvider {
            return Provider()
        },
    })
}

func getSampleResponse(username, password string) (string, error) {
    apiURL := "https://api.example.com"
    authorizationURL := "https://login.example.com"
    client, _ := NewOAuthHTTPClient(authorizationURL, username, password)

    req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2/organizations", apiURL), nil)
    if err != nil {
        return "", err
    }

    response, err := client.Do(req)
    responseBody, err := ioutil.ReadAll(response.Body)

    return string(responseBody), nil
}