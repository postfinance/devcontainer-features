package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var versionRegex *regexp.Regexp = regexp.MustCompile(`(?m)^v(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

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
	installer.HandleGitHubOverride(downloadUrl, "cli/cli", "github-cli-download-url")

	// Create and process the feature
	feature := installer.NewFeature("GitHub CLI", false,
		&ghComponent{
			ComponentBase: installer.NewComponentBase("GitHub CLI", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type ghComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *ghComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("cli", "cli")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *ghComponent) InstallVersion(version *gover.Version) error {
	// Map architecture
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	// Download the file
	fileName := fmt.Sprintf("gh_%s_linux_%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "GitHub CLI"); err != nil {
		return err
	}
	defer os.Remove(fileName)
	// Extract to temp directory
	tempDir, err := os.MkdirTemp("", "github-cli-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, true); err != nil {
		return err
	}
	// Install binary (located at bin/gh inside the extracted folder)
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin(filepath.Join(tempDir, "bin", "gh"), "gh"); err != nil {
		return err
	}
	return nil
}
