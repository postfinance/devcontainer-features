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
	versionsUrl := flag.String("versionsUrl", "", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file (global/per-feature overrides)
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(versionsUrl, "https://releases.jfrog.io/artifactory/jfrog-cli/v2-jf/", "jfrog-cli-versions-url")
	installer.HandleOverride(downloadUrl, "https://releases.jfrog.io/artifactory/jfrog-cli", "jfrog-cli-download-url")

	// Create and process the feature
	feature := installer.NewFeature("JFrog CLI", false,
		&jfrogCliComponent{
			ComponentBase: installer.NewComponentBase("JFrog CLI", *version),
			versionsUrl:   *versionsUrl,
			downloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type jfrogCliComponent struct {
	*installer.ComponentBase
	versionsUrl string
	downloadUrl string
}

func (c *jfrogCliComponent) GetAllVersions() ([]*gover.Version, error) {
	allVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndex(
		c.versionsUrl,
		regexp.MustCompile(`^.*<a href="([0-9][0-9\.]+)/">.*$`),
		versionRegex)
	if err != nil {
		return nil, err
	}
	return allVersions, nil
}

func (c *jfrogCliComponent) InstallVersion(version *gover.Version) error {
	// Map architecture (like in Go installer)
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	// Download the file (filename pattern is the same, just arch changes)
	downloadUrl := fmt.Sprintf("%s/v2-jf/%s/jfrog-cli-linux-%s/jf", c.downloadUrl, version.Raw, archPart)
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
