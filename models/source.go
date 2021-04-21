package models

type Source struct {
	ApiKey   string `json:"api_key"`
	Org      string `json:"org"`
	Project  string `json:"project"`
	Env      string `json:"env"`
	ApiURL   string `json:"api_url"`
}
