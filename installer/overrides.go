package installer

import (
	"fmt"
	"os"
	"strings"
)

// Loads the feature overrides from a specified location.
func LoadOverrides() error {
	var fileContent []byte
	var err error
	if overrideLocation := os.Getenv("DEV_FEATURE_OVERRIDE_LOCATION"); overrideLocation != "" {
		// Load the overrides from the specified location
		if strings.HasPrefix(overrideLocation, "http://") || strings.HasPrefix(overrideLocation, "https://") {
			// Load from URL
			fileContent, err = Tools.Download.AsBytes(overrideLocation)
			if err != nil {
				return fmt.Errorf("error downloading override file: %v", err)
			}
		} else {
			// Load from file
			fileContent, err = os.ReadFile(overrideLocation)
			if err != nil {
				return fmt.Errorf("error reading override file: %v", err)
			}
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
}

func HandleGitHubOverride(downloadUrlBase *string, downloadUrlPath *string, gitHubPath string, key string) {
	if *downloadUrlBase == "" {
		baseKey := keyToEnv(key) + "_BASE"
		if envValue := os.Getenv(baseKey); envValue != "" {
			*downloadUrlBase = envValue
		} else if envValue := os.Getenv("GITHUB_DOWNLOAD_URL_BASE"); envValue != "" {
			*downloadUrlBase = envValue
		} else {
			*downloadUrlBase = "https://github.com"
		}
	}
	if *downloadUrlPath == "" {
		pathKey := keyToEnv(key) + "_PATH"
		if envValue := os.Getenv(pathKey); envValue != "" {
			*downloadUrlPath = envValue
		} else {
			*downloadUrlPath = fmt.Sprintf("%s/releases/download", gitHubPath)
		}
	}
}

func keyToEnv(key string) string {
	return "DEV_FEATURE_OVERRIDE_" + strings.ToUpper(strings.ReplaceAll(key, "-", "_"))
}
