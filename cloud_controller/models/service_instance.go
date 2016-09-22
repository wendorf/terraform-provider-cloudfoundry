package models

type ServiceInstance struct {
	Name            string `json:"name"`
	SpaceGUID       string `json:"space_guid,omitempty"`
	ServicePlanGUID string `json:"service_plan_guid"`
}

type ServiceInstanceWrapper struct {
	Metadata     V2Metadata
	ServiceInstance ServiceInstance `json:"entity"`
}