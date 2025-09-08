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
	version := flag.String("version", "latest", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Apply override logic for URLs
	installer.HandleOverride(downloadUrl, "https://github.com", "goreleaser-download-url")

	// Create and process the feature
	feature := installer.NewFeature("GoReleaser", true,
		&goreleaserComponent{
			ComponentBase: installer.NewComponentBase("GoReleaser", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type goreleaserComponent struct {
	*installer.ComponentBase
	DownloadUrl string
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
	downloadUrl := fmt.Sprintf("%s/goreleaser/goreleaser/releases/download/%s/goreleaser_Linux_x86_64.tar.gz", c.DownloadUrl, versionTag)
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
