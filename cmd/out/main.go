package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"

	"github.com/cycloidio/cycloid-resource/models"
	"github.com/cycloidio/cycloid-resource/helpers"
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

// terracost will run a terracost estimation
func terracost(org, tfplan, apiURL string) ([]byte, models.GenericVersion, []models.Metadata, error) {
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
		errS := err.Error()
		// Get CLI stderr in case of error
		if ee, ok := err.(*exec.ExitError); ok {
			errS = string(ee.Stderr)
		}
		return nil, nil, nil, fmt.Errorf("unable to estimate terraform costs: %v, %s\n", out, errS)
	}

	// Output the terracost estimate JSON that will be used by the cycloid console
	fmt.Fprintln(os.Stderr, string(out))

	var res Estimation
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal from cy output: %w\n", err)
	}

	var version models.TerraCostVersion
	version.BuildID = os.Getenv("BUILD_ID")
	version.PlannedCost = res.PlannedCost
	version.PriorCost = res.PriorCost

	metadatas := []models.Metadata{
		models.Metadata{Name: "planned_cost", Value: res.PlannedCost},
		models.Metadata{Name: "prior_cost", Value: res.PriorCost},
	}

	return out, version, metadatas, nil
}

// terracost will run an infrapolicy check
func infrapolicy(org, project, env, tfplan, apiURL string) ([]byte, models.GenericVersion, []models.Metadata, error) {
	cmdArgs := []string{
		"infrapolicy",
		"validate",
		"--org",
		org,
		"--api-url",
		apiURL,
		"--env",
		env,
		"--project",
		project,
		"--plan-path",
		tfplan,
		"-o",
		"json",
	}
	out, err := exec.Command("cy", cmdArgs...).Output()
	if err != nil {
		errS := err.Error()
		// Get CLI stderr in case of error
		if ee, ok := err.(*exec.ExitError); ok {
			errS = string(ee.Stderr)
		}
		return nil, nil, nil, fmt.Errorf("unable to estimate infrapolicy plan: %v, %s\n", out, errS)
	}

	var res Result
	if err := json.Unmarshal(out, &res); err != nil {
		return nil, nil, nil, fmt.Errorf("unable to unmarshal from cy output: %w\n", err)
	}

	var (
		version   models.InfraPolicyVersion
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

	return out, version, metadatas, nil

}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "Expected path to sources as first argument")
		os.Exit(1)
	}

	sourceDir := os.Args[1]
	if err := os.Chdir(sourceDir); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to access source dir: %v", err)
		os.Exit(1)
	}

	var req models.OutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to read from stdin: %v", err)
		os.Exit(1)
	}

	if err := req.Source.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Resource configuration error: %v", err)
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

	if out, err := exec.Command("cy", loginArgs...).Output(); err != nil {
		errS := err.Error()
		// Get CLI stderr in case of error
		if ee, ok := err.(*exec.ExitError); ok {
			errS = string(ee.Stderr)
		}
		fmt.Fprintf(os.Stderr, "Unable to login to Cycloid: %v, %s\n", out, errS)
		os.Exit(1)
	}

	var (
		version      models.GenericVersion
		metadatas    []models.Metadata
		err          error
		cyJSONOutput []byte
	)

	switch feature, _ := req.Source.GetFeature(); feature {
	case models.TerraCost:
		cyJSONOutput, version, metadatas, err = terracost(req.Source.Org, req.Params.TFPlanPath, req.Source.ApiURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to run terracost check: %v", err)
			os.Exit(1)
		}

	case models.InfraPolicy:
		cyJSONOutput, version, metadatas, err = infrapolicy(req.Source.Org, req.Source.Project, req.Source.Env, req.Params.TFPlanPath, req.Source.ApiURL)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to run infrapolicy: %v", err)
			os.Exit(1)
		}

	default:
		fmt.Fprintf(os.Stderr, "Unknow configured feature named : %s", feature)
		os.Exit(2)
	}

	resp := &models.OutResponse{
		Version:   version,
		Metadatas: metadatas,
	}

	output, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal to output: %v", err)
		os.Exit(1)
	}

	// Write outputs files
	cyJSONFilePath := path.Join(sourceDir, "cy-output.json")
	outputFilePath := path.Join(sourceDir, "output.json")

	if err := helpers.WriteInFile(cyJSONFilePath, string(cyJSONOutput)); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write cy-output.json output file: %v", err)
		os.Exit(1)
	}
	if err := helpers.WriteInFile(outputFilePath, string(output)); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to write output.json output file: %v", err)
		os.Exit(1)
	}

	// Print result
	fmt.Println(string(output))
}
