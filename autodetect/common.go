package autodetect

import (
	"errors"
	"fmt"
	"os"
)

// WriteOrAppendToFile writes the content to the specified file. If the file exists, it appends the content.
func WriteOrAppendToFile(filePath, content string) error {
	if _, err := os.Stat(filePath); errors.Is(err, os.ErrNotExist) {
		f, err := os.Create(filePath)
		if err != nil {
			return fmt.Errorf("Error creating file %s: %w", filePath, err)
		}
		defer f.Close()
		_, err = f.WriteString(content)
		if err != nil {
			return fmt.Errorf("Error writing to file %s: %w", filePath, err)
		}
	} else {
		f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("Error opening file %s: %w", filePath, err)
		}
		defer f.Close()
		_, err = f.WriteString(content)
		if err != nil {
			return fmt.Errorf("Error writing to file %s: %w", filePath, err)
		}
	}
	return nil
}
