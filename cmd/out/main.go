package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"github.com/cycloidio/infrapolicy-resource/models"
)

type Result struct {
	Advisories []Check `json:"advisories"`
	Warnings   []Check `json:"warnings"`
	Criticals  []Check `json:"criticals"`
}

type Check struct {
	Reasons []string `json:"reasons"`
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected path to sources as first argument")
		os.Exit(1)
	}
	sourceDir := os.Args[1]
	if err := os.Chdir(sourceDir); err != nil {
		fmt.Printf("unable to access source dir: %v\n", err)
		os.Exit(1)
	}

	var req models.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Printf("unable to read from stdin: %v\n", err)
		os.Exit(1)
	}

	if req.Source.Email == "" || req.Source.Password == "" {
		fmt.Printf("email and password are required")
		os.Exit(1)
	}

	if req.Source.Org == "" || req.Source.Env == "" || req.Source.Project == "" {
		fmt.Printf("org, env and project are required")
		os.Exit(1)
	}

	loginArgs := []string{
		"login",
		"--org",
		req.Source.Org,
		"--email",
		req.Source.Email,
		"--password",
		req.Source.Password,
		"--api-url",
		req.Source.ApiURL,
	}

	if _, err := exec.Command("cy", loginArgs...).Output(); err != nil {
		fmt.Printf("unable to login to Cycloid: %v\n", err)
		os.Exit(1)
	}

	validateArgs := []string{
		"infrapolicy",
		"validate",
		"--org",
		req.Source.Org,
		"--api-url",
		req.Source.ApiURL,
		"--env",
		req.Source.Env,
		"--project",
		req.Source.Project,
		"--plan-path",
		req.Params.TFPlanPath,
		"-o",
		"json",
	}

	out, err := exec.Command("cy", validateArgs...).Output()
	if err != nil {
		fmt.Printf("unable to validate terraform plan: %v\n", err)
		os.Exit(1)
	}

	var res Result
	if err := json.Unmarshal(out, &res); err != nil {
		fmt.Printf("unable to unmarshal from cy output: %v\n", err)
		os.Exit(1)
	}

	var (
		version   models.Version
		metadatas []models.Metadata
	)

	if len(res.Criticals) > 0 {
		version.Criticals = strconv.Itoa(len(res.Criticals))
		m := models.Metadata{Name: "criticals"}
		for _, critical := range res.Criticals {
			for _, check := range critical.Reasons {
				m.Value += fmt.Sprintf("%s\n", check)
			}
		}
		metadatas = append(metadatas, m)
	} else {
		version.Criticals = "0"
	}
	if len(res.Warnings) > 0 {
		version.Warnings = strconv.Itoa(len(res.Warnings))
		m := models.Metadata{Name: "warnings"}
		for _, warning := range res.Warnings {
			for _, check := range warning.Reasons {
				m.Value += fmt.Sprintf("%s\n", check)
			}
		}
		metadatas = append(metadatas, m)
	} else {
		version.Warnings = "0"
	}
	if len(res.Advisories) > 0 {
		version.Advisories = strconv.Itoa(len(res.Advisories))
		m := models.Metadata{Name: "advisories"}
		for _, advisory := range res.Advisories {
			for _, check := range advisory.Reasons {
				m.Value += fmt.Sprintf("%s\n", check)
			}
		}
		metadatas = append(metadatas, m)
	} else {
		version.Advisories = "0"
	}

	version.BuildID = os.Getenv("BUILD_ID")
	resp := models.OutResponse{
		Version:   version,
		Metadatas: metadatas,
	}

	output, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("unable to marshal to output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
