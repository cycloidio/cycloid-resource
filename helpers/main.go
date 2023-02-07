package helpers


import (
	"io"
	"os"
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
