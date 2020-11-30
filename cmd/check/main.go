package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/cycloidio/infrapolicy-resource/models"
)

func main() {
	var req models.InRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		log.Fatalf("Failed to read InRequest: %s", err)
	}

	out := []models.Version{
		req.Version,
	}

	output, err := json.Marshal(out)
	if err != nil {
		fmt.Printf("unable to marshal to output: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(output))

}
