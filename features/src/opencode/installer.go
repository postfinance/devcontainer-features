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
	installer.HandleGitHubOverride(downloadUrl, "anomalyco/opencode", "opencode-download-url")

	// Create and process the feature
	feature := installer.NewFeature("opencode", false,
		&opencodeComponent{
			ComponentBase: installer.NewComponentBase("opencode", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type opencodeComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *opencodeComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("anomalyco", "opencode")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *opencodeComponent) InstallVersion(version *gover.Version) error {
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
		fileName = fmt.Sprintf("opencode-linux-%s-musl.tar.gz", archPart)
	} else {
		fileName = fmt.Sprintf("opencode-linux-%s.tar.gz", archPart)
	}

	// Download the file
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "opencode"); err != nil {
		return err
	}

	// Extract to a temp directory
	tempDir, err := os.MkdirTemp("", "opencode-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, false); err != nil {
		return err
	}

	// Move the binary to /usr/local/bin/opencode
	if err := installer.Tools.FileSystem.MoveFile(filepath.Join(tempDir, "opencode"), "/usr/local/bin/opencode"); err != nil {
		return err
	}

	// Apply executable permissions
	if err := os.Chmod("/usr/local/bin/opencode", 0755); err != nil {
		return err
	}

	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}

	return nil
}
