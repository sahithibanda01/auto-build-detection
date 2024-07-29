package autodetect

import (
	"fmt"
	"errors"
	"os"
	"path/filepath"
)

type gradleInjecter struct{}

func newGradleInjecter() *gradleInjecter {
	return &gradleInjecter{}
}

func (*gradleInjecter) InjectTool() error {
	accountID := os.Getenv("HARNESS_ACCOUNT_ID")
	token := os.Getenv("HARNESS_PAT")
	endpoint := os.Getenv("HARNESS_CACHE_SERVICE_ENDPOINT")

	// Check if environment variables are set
	if accountID == "" || token == "" || endpoint == "" {
		return errors.New("please set HARNESS_ACCOUNT_ID, HARNESS_PAT, and HARNESS_CACHE_SERVICE_ENDPOINT")
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	gradleDir := filepath.Join(homeDir, ".gradle")
	initGradleFile := filepath.Join(gradleDir, "init.gradle")
	gradlePropertiesFile := filepath.Join(gradleDir, "gradle.properties")

	// Ensure the ~/.gradle directory exists
	err = os.MkdirAll(gradleDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating .gradle directory:", err)
	}

	// Define the content to be written to init.gradle
	initGradleContent := fmt.Sprintf(`
initscript {
    repositories {
        mavenCentral()
    }
    dependencies {
        classpath 'io.harness:gradle-cache:0.0.2'
    }
}
// Apply the plugin to the Settings object
gradle.settingsEvaluated { settings ->
    settings.pluginManager.apply(io.harness.HarnessBuildCache)
    settings.buildCache {
            local {
                enabled = false
            }
            remote(io.harness.Cache) {
                accountId = '%s'
                token = '%s'
                push = "true"
                endpoint = '%s'
            }
        }
}
`, accountID, token, endpoint)

	// Write or append the content to the init.gradle file
	err = WriteOrAppendToFile(initGradleFile, initGradleContent)
	if err != nil {
		return fmt.Errorf("error writing to init.gradle file: %w", err)
	}

	// Write or append the content to the gradle.properties file

	gradlePropertiesContent := "org.gradle.caching=true\n"

	err = WriteOrAppendToFile(gradlePropertiesFile, gradlePropertiesContent)
	if err != nil {
		return fmt.Errorf("error writing to gradle.properties file: %w", err)
	}

	return nil
}
