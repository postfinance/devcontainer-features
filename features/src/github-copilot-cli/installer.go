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
	installer.HandleGitHubOverride(downloadUrl, "github/copilot-cli", "github-copilot-cli-download-url")

	// Create and process the feature
	feature := installer.NewFeature("GitHub Copilot CLI", false,
		&githubCopilotCliComponent{
			ComponentBase: installer.NewComponentBase("GitHub Copilot CLI", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type githubCopilotCliComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *githubCopilotCliComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("github", "copilot-cli")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *githubCopilotCliComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}

	// Download the archive
	fileName := fmt.Sprintf("copilot-linux-%s.tar.gz", archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "GitHub Copilot CLI"); err != nil {
		return err
	}
	defer os.Remove(fileName)

	// Extract to a temp directory
	tempDir, err := os.MkdirTemp("", "github-copilot-cli-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, false); err != nil {
		return err
	}

	// Install the binary
	if err := installer.Tools.System.InstallBinaryToUsrLocalBin(filepath.Join(tempDir, "copilot"), "copilot"); err != nil {
		return err
	}

	return nil
}
