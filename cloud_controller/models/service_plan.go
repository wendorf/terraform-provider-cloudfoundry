package models

type ServicePlan struct {
	Name            string `json:"name"`
}

type ServicePlanCollection struct {
	ServicePlans []ServicePlanWrapper `json:"resources"`
}

type ServicePlanWrapper struct {
	Metadata     V2Metadata
	ServicePlan ServicePlan `json:"entity"`
}