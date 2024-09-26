package models

type Params struct {
	// Terracost
	// TFPlanPath is the path to the terraform plan file
	TFPlanPath string `json:"tfplan_path"`

	// Event
	Title        string            `json:"title"`
	Message      string            `json:"message"`
	YamlVarsFile string            `json:"yaml_vars_file"`
	MessageFile  string            `json:"message_file"`
	Severity     string            `json:"severity"`
	Type         string            `json:"type"`
	Icon         string            `json:"icon"`
	Tags         map[string]string `json:"tags"`
}
