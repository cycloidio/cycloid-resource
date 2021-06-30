package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cycloidio/cycloid-resource/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, "expected output path as first arg")
		os.Exit(1)
	}

	var req models.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "unable to read stdin: %v", err)
		os.Exit(1)
	}

	if err := req.Source.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Resource configuration error: %v", err)
		os.Exit(1)
	}

	feature, _ := req.Source.GetFeature()
	if feature == models.InfraPolicy && req.Version != nil {
		var infraPolicyVersion models.InfraPolicyVersion
		// Now that we know the feature used, we can put back the version from stdin buffer
		// into json to decode it with the appropriated structure
		b, _ := json.Marshal(req.Version)
		if err := json.Unmarshal(b, &infraPolicyVersion); err != nil {
			fmt.Fprintf(os.Stderr, "Unable to decode the version passed as argument: %v", err)
			os.Exit(1)
		}

		criticals, err := strconv.Atoi(infraPolicyVersion.Criticals)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to get number of criticals check: %v", err)
			os.Exit(1)
		}
		warnings, err := strconv.Atoi(infraPolicyVersion.Warnings)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to get number of warnings check: %v", err)
			os.Exit(1)
		}
		if criticals > 0 || warnings > 0 {
			fmt.Fprint(os.Stderr, "critical or warning checks are present, check metadata of your resource for more information")
			os.Exit(1)
		}
	}

	resp := models.OutResponse{
		Version:   req.Version,
		Metadatas: []models.Metadata{},
	}

	output, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal to output: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
