package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSpace() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpaceCreate,
		Read:   resourceSpaceRead,
		Update: resourceSpaceUpdate,
		Delete: resourceSpaceDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceSpaceCreate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSpaceRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSpaceUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceSpaceDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}
