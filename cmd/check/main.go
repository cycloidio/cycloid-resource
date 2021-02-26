package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/cycloidio/infrapolicy-resource/models"
)

func main() {
	var req models.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		fmt.Fprintf(os.Stderr, "unable to decode request: %v", err)
		os.Exit(1)
	}

	out := []models.Version{
		req.Version,
	}

	output, err := json.Marshal(out)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to marshal output: %v", err)
		os.Exit(1)
	}
	fmt.Println(string(output))

}
