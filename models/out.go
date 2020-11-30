package models

type OutRequest struct {
	Source Source `json:"source"`
	Params Params `json:"params"`
}

type OutResponse struct {
	Version   Version    `json:"version"`
	Metadatas []Metadata `json:"metadata"`
}

type Metadata struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
