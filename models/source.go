package models

type Source struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Org      string `json:"org"`
	Project  string `json:"project"`
	Env      string `json:"env"`
	ApiURL   string `json:"api_url"`
}
