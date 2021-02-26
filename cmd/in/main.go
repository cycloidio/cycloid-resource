package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/cycloidio/infrapolicy-resource/models"
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

	criticals, err := strconv.Atoi(req.Version.Criticals)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get number of criticals check: %v", err)
		os.Exit(1)
	}
	warnings, err := strconv.Atoi(req.Version.Warnings)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to get number of warnings check: %v", err)
		os.Exit(1)
	}

	if criticals > 0 || warnings > 0 {
		fmt.Fprint(os.Stderr, "critical or warning checks are present, check metadata of your resource for more information")
		os.Exit(1)
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
