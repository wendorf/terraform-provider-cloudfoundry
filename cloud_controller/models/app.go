package models

type App struct {
	GUID          string `json:"guid,omitempty"`
	Name          string `json:"name"`
	Relationships AppRelationships `json:"relationships"`
}

type AppRelationships struct {
	Space GUID `json:"space"`
}

type GUID struct {
	GUID string `json:"guid"`
}