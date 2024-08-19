package autodetect

import (
	"errors"
	"fmt"
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
	gradlePluginVersion := "0.0.2" // make this configurable via command line input or env variable
	gradleCachePush := "true"

	// Check if environment variables are set
	if accountID == "" || token == "" || endpoint == "" {
		return errors.New("please set HARNESS_ACCOUNT_ID, HARNESS_PAT, and HARNESS_CACHE_SERVICE_ENDPOINT")
	}

	gradlePropertiesContent := "org.gradle.caching=true\n"
	// Define the content to be written to init.gradle
	initGradleContent := fmt.Sprintf(`
initscript {
    repositories {
        mavenCentral()
    }
    dependencies {
        classpath 'io.harness:gradle-cache:%s'
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
                accountId = System.getenv('HARNESS_ACCOUNT_ID')
                token = System.getenv('HARNESS_PAT')
                push = "%s"
                endpoint = System.getenv('HARNESS_CACHE_SERVICE_ENDPOINT')
            }
        }
}
`, gradlePluginVersion, gradleCachePush)
	// Injecting the above configs in gradle files
	// For $GRADLE_HOME
	gradleHome := os.Getenv("GRADLE_HOME")
	if gradleHome != "" {
		// $GRADLE_HOME/init.d/init.gradle file
		gradleHomeInit := filepath.Join(gradleHome, "init.d")
		err := os.MkdirAll(gradleHomeInit, os.ModePerm)
		initGradleHomeFile := filepath.Join(gradleHomeInit, "init.gradle")
		err = WriteOrAppendToFile(initGradleHomeFile, initGradleContent)
		if err != nil {
			return fmt.Errorf("error writing to $GRADLE_HOME/init.d/init.gradle file: %w", err)
		}
		// $GRADLE_HOME/init.d/gradle.properties file
		gradleHomePropertiesFile := filepath.Join(gradleHome, "gradle.properties")
		err = WriteOrAppendToFile(gradleHomePropertiesFile, gradlePropertiesContent)
		if err != nil {
			return fmt.Errorf("error writing to $GRADLE_HOME/init.d/gradle.properties file: %w", err)
		}
	}

	// For $GRADLE_USER_HOME
	gradleUserHome := os.Getenv("GRADLE_USER_HOME")
	if gradleUserHome != "" {
		// $GRADLE_USER_HOME/init.d/init.gradle file
		gradleUserHomeInit := filepath.Join(gradleUserHome, "init.d")
		err := os.MkdirAll(gradleUserHomeInit, os.ModePerm)
		initGradleUserHomeFile := filepath.Join(gradleUserHomeInit, "init.gradle")
		err = WriteOrAppendToFile(initGradleUserHomeFile, initGradleContent)
		if err != nil {
			return fmt.Errorf("error writing to $GRADLE_USER_HOME/init.d/init.gradle file: %w", err)
		}

		// $GRADLE_USER_HOME/gradle.properties file
		gradleUserHomePropertiesFile := filepath.Join(gradleUserHome, "gradle.properties")
		err = WriteOrAppendToFile(gradleUserHomePropertiesFile, gradlePropertiesContent)
		if err != nil {
			return fmt.Errorf("error writing to $GRADLE_USER_HOME/gradle.properties file: %w", err)
		}
	}

	// For ~/.gradle
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("error getting user home directory: %w", err)
	}

	gradleDir := filepath.Join(homeDir, ".gradle")
	gradleInitdDir := filepath.Join(gradleDir, "init.d")
	initGradleFile := filepath.Join(gradleInitdDir, "init.gradle")
	gradlePropertiesFile := filepath.Join(gradleInitdDir, "gradle.properties")

	// Create the ~/.gradle/init.d directory if it does not exist
	err = os.MkdirAll(gradleInitdDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("error creating ~/.gradle/init.d directory: %s", err)
	}

	// ~/.gradle/init.d/init.gradle file
	err = WriteOrAppendToFile(initGradleFile, initGradleContent)
	if err != nil {
		return fmt.Errorf("error writing to ~/.gradle/init.d/init.gradle file: %w", err)
	}

	// ~/.gradle/gradle.properties file
	err = WriteOrAppendToFile(gradlePropertiesFile, gradlePropertiesContent)
	if err != nil {
		return fmt.Errorf("error writing to ~/.gradle/gradle.properties file: %w", err)
	}

	return nil
}
