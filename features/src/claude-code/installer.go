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
	installer.HandleGitHubOverride(downloadUrl, "anthropics/claude-code", "claude-code-download-url")

	// Create and process the feature
	feature := installer.NewFeature("claude-code", false,
		&claudeCodeComponent{
			ComponentBase: installer.NewComponentBase("claude-code", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type claudeCodeComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *claudeCodeComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("anthropics", "claude-code")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *claudeCodeComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}

	// Use musl variant for Alpine (musl libc)
	osInfo, err := installer.Tools.System.GetOsInfo()
	if err != nil {
		return err
	}
	var fileName string
	if osInfo.IsAlpine() {
		fileName = fmt.Sprintf("claude-linux-%s-musl.tar.gz", archPart)
	} else {
		fileName = fmt.Sprintf("claude-linux-%s.tar.gz", archPart)
	}

	// Download the file from GitHub releases
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "claude-code"); err != nil {
		return err
	}
	defer os.Remove(fileName)

	// Extract to a temp directory
	tempDir, err := os.MkdirTemp("", "claude-code-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, false); err != nil {
		return err
	}

	// Install the binary
	return installer.Tools.System.InstallBinaryToUsrLocalBin(filepath.Join(tempDir, "claude"), "claude")
}
