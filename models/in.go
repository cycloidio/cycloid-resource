package models

type InRequest struct {
	Source  Source         `json:"source"`
	Version GenericVersion `json:"version"`
	Params  Params         `json:"params"`
}
