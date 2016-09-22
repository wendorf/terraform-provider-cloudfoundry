package models

type Route struct {
	Host       string `json:"host"`
	DomainGUID string `json:"domain_guid"`
	SpaceGUID  string `json:"space_guid"`
}