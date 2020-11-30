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
		fmt.Println("expected output path as first arg")
		os.Exit(1)
	}

	var req models.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Println("unable to read stdin: %v\n", err)
		os.Exit(1)
	}

	criticals, err := strconv.Atoi(req.Version.Criticals)
	if err != nil {
		fmt.Println("unable to get number of criticals check")
		os.Exit(1)
	}
	warnings, err := strconv.Atoi(req.Version.Warnings)
	if err != nil {
		fmt.Println("unable to get number of criticals check")
		os.Exit(1)
	}

	if criticals > 0 || warnings > 0 {
		fmt.Println("critical or warning checks are present, check metadata of your resource for more information")
		os.Exit(1)
	}

	resp := models.OutResponse{
		Version:   req.Version,
		Metadatas: []models.Metadata{},
	}

	output, err := json.Marshal(resp)
	if err != nil {
		fmt.Printf("unable to marshal to output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
