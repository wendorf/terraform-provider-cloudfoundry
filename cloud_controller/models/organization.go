package models

type Organization struct {
	Name string `json:"name"`
}

type OrganizationCollection struct {
	Organizations []OrganizationWrapper `json:"resources"`
}

type OrganizationWrapper struct {
	Metadata     V2Metadata
	Organization Organization `json:"entity"`
}