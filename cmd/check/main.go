package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cycloidio/cycloid-resource/models"
)

func main() {
	var req models.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode request: %v", err)
		os.Exit(1)
	}

	if err := req.Source.Validate(); err != nil {
		fmt.Fprintf(os.Stderr, "Resource configuration error: %v", err)
		os.Exit(1)
	}

	out := []models.GenericVersion{
		req.Version,
	}

	output, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal output: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(output))
}
