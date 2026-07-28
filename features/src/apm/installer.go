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
	installer.HandleGitHubOverride(downloadUrl, "microsoft/apm", "apm-download-url")

	// Create and process the feature
	feature := installer.NewFeature("APM", true,
		&apmComponent{
			ComponentBase: installer.NewComponentBase("APM", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type apmComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *apmComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("microsoft", "apm")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *apmComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}

	// Download the archive
	archiveName := fmt.Sprintf("apm-linux-%s", archPart)
	fileName := archiveName + ".tar.gz"
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "APM"); err != nil {
		return err
	}
	defer os.Remove(fileName)

	// Extract to a temp directory
	tempDir, err := os.MkdirTemp("", "apm-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, false); err != nil {
		return err
	}

	// Install the binary (located in apm-linux-{arch}/apm inside the archive)
	return installer.Tools.System.InstallBinaryToUsrLocalBin(filepath.Join(tempDir, archiveName, "apm"), "apm")
}
