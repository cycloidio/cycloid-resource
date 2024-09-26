package models

type GenericVersion interface{}

type InfraPolicyVersion struct {
	BuildID    string `json:"build_id"`
	Criticals  string `json:"criticals"`
	Warnings   string `json:"warnings"`
	Advisories string `json:"advisories"`
}

type TerraCostVersion struct {
	BuildID     string `json:"build_id"`
	PlannedCost string `json:"planned_cost"`
	PriorCost   string `json:"prior_cost"`
}

type EventVersion struct {
	BuildID string `json:"build_id"`
}
