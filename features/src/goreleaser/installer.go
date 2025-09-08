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

var githubReleasesBaseUrl = "https://github.com"
var versionRegex *regexp.Regexp = regexp.MustCompile(`^v(\d+)\.(\d+)\.(\d+)?$`)

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
	version := flag.String("version", "latest", "The version of GoReleaser to install.")
	flag.Parse()

	// Create and process the feature
	feature := installer.NewFeature("GoReleaser", true,
		&goreleaserComponent{
			ComponentBase: installer.NewComponentBase("GoReleaser", *version),
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type goreleaserComponent struct {
	*installer.ComponentBase
}

func (c *goreleaserComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("goreleaser", "goreleaser")
	if err != nil {
		return nil, err
	}
	return installer.Tools.GitHub.ParseVersionFromTags(tags, versionRegex, "nightly")
}

func (c *goreleaserComponent) InstallVersion(version *gover.Version) error {
	// Prepend "v" to version.Raw if it does not already start with "v"
	// This can happen if the vesion is not resolved and the input is without v
	versionTag := version.Raw
	if len(versionTag) > 0 && versionTag[0] != 'v' {
		versionTag = "v" + versionTag
	}
	// Download the file
	downloadUrl := fmt.Sprintf("%s/goreleaser/goreleaser/releases/download/%s/goreleaser_Linux_x86_64.tar.gz", githubReleasesBaseUrl, versionTag)
	fileName := "goreleaser.tar.gz"
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "goreleaser"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "/usr/local/goreleaser", false); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}
