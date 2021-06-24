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

type Estimation struct {
	PlannedCost string `json:"planned_cost"`
	PriorCost   string `json:"prior_cost"`
}

// terracost will run a terracost
// estimation
func terracost(org, tfplan, apiURL string) ([]models.Metadata, error) {
	terracostArgs := []string{
		"terracost",
		"estimate",
		"--org",
		org,
		"--plan-path",
		tfplan,
		"--api-url",
		apiURL,
		"-o",
		"json",
	}
	out, err := exec.Command("cy", terracostArgs...).Output()
	if err != nil {
		return nil, fmt.Errorf("unable to estimate terraform plan costs: %w\n", err)
	}

	// Output the terracost estimate JSON that will be used by the cycloid console
	fmt.Fprintln(os.Stderr, string(out))

	var res Estimation
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, fmt.Errorf("unable to unmarshal from cy output: %w\n", err)
	}

	return []models.Metadata{
		models.Metadata{Name: "planned_cost", Value: res.PlannedCost},
		models.Metadata{Name: "prior_cost", Value: res.PriorCost},
	}, nil
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "expected path to sources as first argument")
		os.Exit(1)
	}
	sourceDir := os.Args[1]
	if err := os.Chdir(sourceDir); err != nil {
		fmt.Fprintf(os.Stderr, "unable to access source dir: %v", err)
		os.Exit(1)
	}

	var req models.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "unable to read from stdin: %v", err)
		os.Exit(1)
	}

	if req.Source.ApiKey == "" {
		fmt.Fprint(os.Stderr, "api_key is required")
		os.Exit(1)
	}

	if req.Source.Org == "" || req.Source.Env == "" || req.Source.Project == "" {
		fmt.Fprint(os.Stderr, "org, env and project are required")
		os.Exit(1)
	}

	loginArgs := []string{
		"login",
		"--org",
		req.Source.Org,
		"--api-key",
		req.Source.ApiKey,
		"--api-url",
		req.Source.ApiURL,
	}

	if _, err := exec.Command("cy", loginArgs...).Output(); err != nil {
		fmt.Fprintf(os.Stderr, "unable to login to Cycloid: %v", err)
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
		fmt.Fprintf(os.Stderr, "unable to validate terraform plan: %v", err)
		os.Exit(1)
	}

	var res Result
	if err := json.Unmarshal(out, &res); err != nil {
		fmt.Fprintf(os.Stderr, "unable to unmarshal from cy output: %v", err)
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

	if req.Params.Terracost {
		estimations, err := terracost(req.Source.Org, req.Params.TFPlanPath, req.Source.ApiURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to run terracost check: %v", err)
			os.Exit(1)
		}
		metadatas = append(metadatas, estimations...)
	}

	version.BuildID = os.Getenv("BUILD_ID")
	resp := models.OutResponse{
		Version:   version,
		Metadatas: metadatas,
	}

	output, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal to output: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
