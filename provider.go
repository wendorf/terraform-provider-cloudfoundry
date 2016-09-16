package main

import (
    "github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
    return &schema.Provider{
        ResourcesMap: map[string]*schema.Resource{
            "cloudfoundry_space": resourceSpace(),
        },
    }
}
