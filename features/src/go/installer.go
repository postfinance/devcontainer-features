package main

import (
	"builder/installer"
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

const versionsIndexUrl = "https://go.dev/dl/?mode=json&include=all"
const latestVersionUrl = "https://go.dev/VERSION?m=text"

var versionRegexp *regexp.Regexp = regexp.MustCompile(`(?m:)^go(?P<d1>\d+)(?:\.(?P<d2>\d+))?(?:\.(?P<d3>\d+))?(?:(?P<s4>[a-z]+)(?P<d5>\d+)?)?$`)

//////////
// Main
//////////

func main() {
	if err := runMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	// Handle the flags
	version := flag.String("version", "latest", "The version of Go to install.")
	downloadRegistryBase := flag.String("downloadRegistryBase", "", "The download registry to use for Go binaries.")
	downloadRegistryPath := flag.String("downloadRegistryPath", "", "The download registry path to use for Go binaries.")
	flag.Parse()

	fill(downloadRegistryBase, "https://dl.google.com", "")
	fill(downloadRegistryPath, "/go", "")

	// Create and process the feature
	feature := installer.NewFeature("Go", true,
		&goComponent{
			ComponentBase:        installer.NewComponentBase("Go", *version),
			DownloadRegistryBase: *downloadRegistryBase,
			DownloadRegistryPath: *downloadRegistryPath,
		})
	return feature.Process()
}

func fill(passedValue *string, defaultValue string, key string) {
	if *passedValue == "" {
		// TODO: Get value from env?file?
		// Otherwise set to default value
		*passedValue = defaultValue
	}
}

func buildUrl(base string, parts ...string) (string, error) {
	return url.JoinPath(base, parts...)
}

//////////
// Implementation
//////////

type goComponent struct {
	*installer.ComponentBase
	DownloadRegistryBase string
	DownloadRegistryPath string
}

func (c *goComponent) GetLatestVersion() (*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsString(latestVersionUrl)
	if err != nil {
		return nil, err
	}
	// Only use the first line
	lines := strings.Split(strings.ReplaceAll(versionFileContent, "\r\n", "\n"), "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("no version found in go latest")
	}
	version, err := gover.ParseVersionFromRegex(lines[0], versionRegexp)
	if err != nil {
		return nil, err
	}
	return version, err
}

func (c *goComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(versionsIndexUrl)
	if err != nil {
		return nil, err
	}
	var jsonData []map[string]interface{}
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	versions := []*gover.Version{}
	for _, entry := range jsonData {
		versionString := entry["version"].(string)
		version, err := gover.ParseVersionFromRegex(versionString, versionRegexp)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func (c *goComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	fileName := fmt.Sprintf("%s.linux-amd64.tar.gz", version.Raw)
	downloadUrl, err := buildUrl(c.DownloadRegistryBase, c.DownloadRegistryPath, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Go"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "/usr/local", false); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}
