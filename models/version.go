package models

type Version struct {
	BuildID    string `json:"build_id"`
	Criticals  string `json:"criticals"`
	Warnings   string `json:"warnings"`
	Advisories string `json:"advisories"`
}
