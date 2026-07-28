package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"

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
	// Check Preconditions
	osInfo, err := installer.Tools.System.GetOsInfo()
	if err != nil {
		return fmt.Errorf("failed getting os info: %v", err)
	}
	if osInfo.Vendor == "debian" {
		versionId, err := strconv.Atoi(osInfo.VersionId)
		if err != nil {
			return fmt.Errorf("failed parsing the version number from %s: %v", osInfo.Vendor, err)
		}
		if versionId < 13 {
			return fmt.Errorf("unsupported debian version: %d", versionId)
		}
	}

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

	// Move the extracted archive directory to /opt/apm to preserve the depencencies
	// Can be simplified if the release is changed to a single binary
	installDir := "/opt/apm"
	if err := os.RemoveAll(installDir); err != nil && !os.IsNotExist(err) {
		return err
	}
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return err
	}
	extractedDir := filepath.Join(tempDir, archiveName)
	if err := installer.Tools.FileSystem.MoveFolder(extractedDir, installDir, false); err != nil {
		return err
	}

	// Create a symlink from /usr/local/bin/apm to the installed binary
	binaryPath := filepath.Join(installDir, "apm")
	return installer.Tools.FileSystem.CreateSymLink(binaryPath, "/usr/local/bin/apm", false)
}
