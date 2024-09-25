package autodetect

import (
	"fmt"
	"os"
)

// WriteOrAppendToFile writes the content to the specified file. If the file exists, it appends the content.
func WriteOrAppendToFile(filePath, content string) error {
	// create or open the file
	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening file %s file: %s\n", filePath, err.Error())
		return fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		fmt.Printf(" writing to file %s file: %s\n", filePath, err.Error())
		return fmt.Errorf("error writing to file %s: %w", filePath, err)
	}
	return nil
}
