package models

type Service struct {
}

type ServiceCollection struct {
	Services []ServiceWrapper `json:"resources"`
}

type ServiceWrapper struct {
	Metadata     V2Metadata
	Service Service `json:"entity"`
}