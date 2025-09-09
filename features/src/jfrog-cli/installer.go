package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

const baseUrl = "https://releases.jfrog.io/artifactory/jfrog-cli"
const allVersionsUrl = "https://releases.jfrog.io/artifactory/jfrog-cli/v2-jf/"

var versionRegex *regexp.Regexp = regexp.MustCompile(`(?m:)^(\d+).(\d+).(\d+)$`)

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
	version := flag.String("version", "latest", "")
	flag.Parse()

	// Create and process the feature
	feature := installer.NewFeature("JFrog CLI", false,
		&jfrogCliComponent{
			ComponentBase: installer.NewComponentBase("JFrog CLI", *version),
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type jfrogCliComponent struct {
	*installer.ComponentBase
}

func (c *jfrogCliComponent) GetAllVersions() ([]*gover.Version, error) {
	allVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndex(
		allVersionsUrl,
		regexp.MustCompile(`^.*<a href="([0-9][0-9\.]+)/">.*$`),
		versionRegex)
	if err != nil {
		return nil, err
	}
	return allVersions, nil
}

func (c *jfrogCliComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	downloadUrl := fmt.Sprintf("%s/v2-jf/%s/jfrog-cli-linux-amd64/jf", baseUrl, version.Raw)
	if err := installer.Tools.Download.ToFile(downloadUrl, "/usr/local/bin/jf", "JF"); err != nil {
		return err
	}
	// Set execute rights
	if err := os.Chmod("/usr/local/bin/jf", 0x755); err != nil {
		return err
	}
	// Add alias to jfrog
	if err := installer.Tools.FileSystem.CreateSymLink("/usr/local/bin/jf", "/usr/local/bin/jfrog", false); err != nil {
		return err
	}
	return nil
}
