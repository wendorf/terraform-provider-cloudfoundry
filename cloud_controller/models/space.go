package models

type Space struct {
	Name             string `json:"name"`
	OrganizationGUID string `json:"organization_guid"`
}

type SpaceCollection struct {
	Spaces []SpaceWrapper `json:"resources"`
}

type SpaceWrapper struct {
	Metadata     V2Metadata
	Space Space `json:"entity"`
}