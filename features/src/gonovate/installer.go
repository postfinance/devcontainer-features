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

var versionRegex *regexp.Regexp = regexp.MustCompile(`(?m:)^v(?P<raw>(\d+).(\d+)\.(\d+))$`)

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
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	// Apply override logic for URLs
	installer.HandleGitHubOverride(downloadUrl, "roemer/gonovate", "gonovate-download-url")

	// Create and process the feature
	feature := installer.NewFeature("Gonovate", true,
		&gonovateComponent{
			ComponentBase: installer.NewComponentBase("Gonovate", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type gonovateComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *gonovateComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("roemer", "gonovate")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *gonovateComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	// Download the file
	fileName := fmt.Sprintf("gonovate-linux-%s-%s.zip", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "gonovate"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractZip(fileName, "/usr/local/bin/", false); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}
