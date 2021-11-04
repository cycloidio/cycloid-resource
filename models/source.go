package models

import (
	"fmt"
	"strings"
)

const InfraPolicy = "infrapolicy"
const TerraCost = "terracost"

type Source struct {
	// Feature is the name of the Cycloid feature eg infrapolicy, terracost
	Feature string `json:"feature"`
	ApiKey  string `json:"api_key"`
	Org     string `json:"org"`
	Project string `json:"project"`
	Env     string `json:"env"`
	ApiURL  string `json:"api_url"`
}

// GetFeature returns the feature configured
func (s Source) GetFeature() (string, error) {
	f := strings.ToLower(s.Feature)

	if f == "" {
		return "", fmt.Errorf("feature field is empty")
	}
	if f != InfraPolicy && f != TerraCost {
		return "", fmt.Errorf("feature field should match %s", strings.Join([]string{InfraPolicy, TerraCost}, ", "))
	}

	return f, nil
}

func (s *Source) Validate() error {
	var err error

	_, ferr := s.GetFeature()
	if ferr != nil {
		err = fmt.Errorf("feature configuration error: %v", ferr)
	}
	if s.ApiKey == "" {
		err = fmt.Errorf("api_key is required")
	}
	if s.Org == "" || s.Env == "" || s.Project == "" {
		err = fmt.Errorf("org, env and project are required")
	}

    // Setting as default our SaaS API as default url
	if s.ApiURL == "" {
		s.ApiURL = "https://http-api.cycloid.io"
	}

	return err
}
