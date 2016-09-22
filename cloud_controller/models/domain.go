package models

type Domain struct {
	Name string `json:"name"`
}

type DomainCollection struct {
	Domains []DomainWrapper `json:"resources"`
}

type DomainWrapper struct {
	Metadata V2Metadata
	Domain   Domain `json:"entity"`
}