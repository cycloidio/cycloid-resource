package models

type Params struct {
	// TFPlanPath is the path to the terraform plan file
	TFPlanPath string `json:"tfplan_path"`
	// Terracost decides to activate or not cost estimation
	Terracost bool `json:"terracost"`
}
