package models

type InRequest struct {
	Source  Source  `json:"source"`
	Version Version `json:"version"`
	Params  Params  `json:"params"`
}
