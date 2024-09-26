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
	bearerToken := os.Getenv("HARNESS_CACHE_SERVICE_TOKEN")
	endpoint := os.Getenv("HARNESS_CACHE_SERVICE_ENDPOINT")
	gradlePluginVersion := os.Getenv("HARNESS_GRADLE_PLUGIN_VERSION")
	gradleCachePush := os.Getenv("HARNESS_CACHE_PUSH")
	localCacheEnabled := os.Getenv("HARNESS_CACHE_LOCAL_ENABLED")

	// Check if environment variables are set
	if accountID == "" || bearerToken == "" || endpoint == "" {
		return errors.New("please set HARNESS_ACCOUNT_ID,HARNESS_CACHE_SERVICE_TOKEN, and HARNESS_CACHE_SERVICE_ENDPOINT")
	}

	// Define the content to be written to gradle.properties
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
                enabled = "%s"
            }
            remote(io.harness.Cache) {
                accountId = System.getenv('HARNESS_ACCOUNT_ID')
                push = "%s"
                endpoint = System.getenv('HARNESS_CACHE_SERVICE_ENDPOINT')
            }
        }
}
`, gradlePluginVersion, localCacheEnabled, gradleCachePush)

	// Injecting the above configs in gradle files
	// For $GRADLE_HOME
	gradleHome := os.Getenv("GRADLE_HOME")
	if gradleHome != "" {
		// $GRADLE_HOME/init.d/init.gradle file
		err := injectGradleFiles(gradleHome, initGradleContent, gradlePropertiesContent)
		if err != nil {
			return errors.New("error in injectiong GRADLE_HOM directory")
		}
	}

	// For $GRADLE_USER_HOME
	gradleUserHome := os.Getenv("GRADLE_USER_HOME")
	if gradleUserHome != "" {
		// $GRADLE_USER_HOME/init.d/init.gradle file
		err := injectGradleFiles(gradleUserHome, initGradleContent, gradlePropertiesContent)
		if err != nil {
			return errors.New("error in injectiong GRADLE_USER_HOME directory")
		}
	}

	// For ~/.gradle
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return errors.New("error getting user home directory: %w")
	}

	gradleDir := filepath.Join(homeDir, ".gradle")
	err = injectGradleFiles(gradleDir, initGradleContent, gradlePropertiesContent)
	if err != nil {
		fmt.Printf("its okay its okay")
		return errors.New(err.Error())
	}
	return nil
}

func injectGradleFiles(gradleDir string, initGradleContent string, gradlePropertiesContent string) error {
	err := os.MkdirAll(gradleDir, os.ModePerm)
	if err != nil {
		fmt.Printf("error creating %s directory: %s\n", gradleDir, err.Error())
		return fmt.Errorf("error creating %s directory: %w", gradleDir, err)
	}
	// $gradleDir/init.gradle file
	gradleHomeInit := filepath.Join(gradleDir, "init.d")
	initGradleHomeFile := filepath.Join(gradleHomeInit, "init.gradle")
	err = WriteOrAppendToFile(initGradleHomeFile, initGradleContent)
	if err != nil {
		fmt.Printf("error writing %s file: %s\n", initGradleContent, err.Error())
		return fmt.Errorf("error writing to %s file: %w", initGradleContent, err)
	}
	// gradleDir/gradle.properties file
	gradleHomePropertiesFile := filepath.Join(gradleDir, "gradle.properties")
	err = WriteOrAppendToFile(gradleHomePropertiesFile, gradlePropertiesContent)
	if err != nil {
		fmt.Printf("error writing %s file: %s\n", gradlePropertiesContent, err.Error())
		return fmt.Errorf("error writing to %s file: %w", gradlePropertiesContent, err)
	}

	return nil
}
