package installer

import (
	"fmt"
	"os"
	"strings"
)

// Loads the feature overrides from a specified location.
func LoadOverrides() error {
	if overrideLocation := os.Getenv("DEV_FEATURE_OVERRIDE_LOCATION"); overrideLocation != "" {
		fileContent, err := ReadFileFromUrlOrLocal(overrideLocation)
		if err != nil {
			return err
		}
		if len(fileContent) > 0 {
			lines := strings.SplitSeq(strings.ReplaceAll(strings.TrimSpace(string(fileContent)), "\r\n", "\n"), "\n")
			for line := range lines {
				if strings.HasPrefix(line, "#") || strings.TrimSpace(line) == "" {
					continue
				}
				parts := strings.SplitN(line, "=", 2)
				if len(parts) == 2 {
					key := strings.TrimSpace(parts[0])
					value := strings.TrimSpace(parts[1])
					if !strings.HasPrefix(key, "DEV_FEATURE_OVERRIDE_") {
						key = fmt.Sprintf("DEV_FEATURE_OVERRIDE_%s", key)
					}
					// Only set if not already defined
					if os.Getenv(key) == "" {
						os.Setenv(key, value)
					}
				}
			}
		}
	}
	return nil
}

// ReadFileFromUrlOrLocal loads a file from a URL or local path and returns its contents as bytes.
func ReadFileFromUrlOrLocal(location string) ([]byte, error) {
	if strings.HasPrefix(location, "http://") || strings.HasPrefix(location, "https://") {
		fileContent, err := Tools.Download.AsBytes(location)
		if err != nil {
			return nil, fmt.Errorf("error downloading file: %v", err)
		}
		return fileContent, nil
	} else {
		fileContent, err := os.ReadFile(location)
		if err != nil {
			return nil, fmt.Errorf("error reading file: %v", err)
		}
		return fileContent, nil
	}
}

func HandleOverride(passedValue *string, defaultValue string, key string) {
	if *passedValue == "" {
		// Convert the key to a compatible format
		key = keyToEnv(key)
		// Try get the value from an environment variable
		if envValue := os.Getenv(key); envValue != "" {
			*passedValue = envValue
			return
		}
		// Otherwise set to default value
		*passedValue = defaultValue
	}
	// Handle "none" value to explicitly unset the value if an override env variable is set to a different value
	// e.g. --configPath="none" and DOCKER_OUT_CONFIG_PATH=https://example.com/config.json
	// If we do not explicitly handle this, you could not unset a value once an override env variable is set
	if strings.ToLower(*passedValue) == "none" {
		*passedValue = ""
	}
}

func HandleGitHubOverride(downloadUrl *string, gitHubPath string, key string) {
	if *downloadUrl == "" {
		if envValue := os.Getenv(keyToEnv(key)); envValue != "" {
			*downloadUrl = envValue
		} else if envValue := os.Getenv(keyToEnv("GITHUB_DOWNLOAD_URL")); envValue != "" {
			*downloadUrl, _ = Tools.Http.BuildUrl(envValue, gitHubPath, "/releases/download")
		} else {
			*downloadUrl, _ = Tools.Http.BuildUrl("https://github.com", gitHubPath, "/releases/download")
		}
	}
}

func keyToEnv(key string) string {
	return "DEV_FEATURE_OVERRIDE_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}
