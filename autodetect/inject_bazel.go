package autodetect

import (
	"fmt"
	"os"
	"path/filepath"
)

type bazelInjecter struct{}

func newBazelInjecter() *bazelInjecter {
	return &bazelInjecter{}
}

func (*bazelInjecter) InjectTool() error {
	homeDir, err := os.UserHomeDir()
	// endpoint := os.Getenv("HARNESS_END_POINT")
	if err != nil {
		return fmt.Errorf("error getting user home directory:", err)
	}

	bazelrcFile := filepath.Join(homeDir, ".bazelrc")

	bazelrcContent := `build --remote_cache=http://localhost:8082/cache/bazel`

	err = WriteOrAppendToFile(bazelrcFile, bazelrcContent)
	if err != nil {
		return fmt.Errorf("error writing to bazelrc file: %w", err)
	}
	return nil
}
