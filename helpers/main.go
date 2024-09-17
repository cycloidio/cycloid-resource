package helpers

import (
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// WriteInFile write content string in a specified file
func WriteInFile(filePath string, content string) error {
	outputFile, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	if _, err := io.WriteString(outputFile, content); err != nil {
		return err
	}
	if err := outputFile.Close(); err != nil {
		return err
	}
	return nil
}

// ReplaceVariables Render bash variables in a string
func ReplaceVariables(input string, vars map[string]string) string {
	// Regex to catch bash variables $FOO and ${FOO}
	re := regexp.MustCompile(`\$\{?([A-Za-z_][A-Za-z0-9_]*)\}?`)

	result := re.ReplaceAllStringFunc(input, func(varStr string) string {
		varName := strings.Trim(varStr, "${}")

		// Get value for the var, if not found lookup in environment variables
		if value, exists := vars[varName]; exists {
			return value
		} else if value, exists := os.LookupEnv(varName); exists {
			return value
		}

		return varStr
	})

	return result
}

// LoadYAMLToMap load Yaml string map
func LoadYAMLToMap(filename string) (map[string]string, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)

	err = yaml.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadFileToString Function to read the content of a text file and return it as a string
func ReadFileToString(filename string) (string, error) {
	// Read the file content
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", err
	}

	// Convert the byte slice to a string and return
	return string(data), nil
}

func MakeStringShellSafe(input string) string {
	// Escape double quotes, backslashes, dollar signs, backticks, and newlines
	replacer := strings.NewReplacer(
		//`"`, `\"`, // Escape double quotes
		`\`, `\\`, // Escape backslashes
		//"$", `\$`,   // Escape dollar signs
		"`", "\\`", // Escape backticks
		"\n", `<br />`, // Escape newlines for safe shell usage
		//"\n", `\n`, // Escape newlines for safe shell usage
	)

	return replacer.Replace(input)
}
