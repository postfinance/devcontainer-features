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

var versionRegexp *regexp.Regexp = regexp.MustCompile(`(\d+)\.(\d+)\.(\d+)`)

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

	installer.HandleOverride(downloadUrl, "https://gitlab.com/gitlab-org/cli/-/releases", "gitlab-cli-download-url")

	// Create and process the feature
	feature := installer.NewFeature("GitLab CLI", false,
		&glabComponent{
			ComponentBase: installer.NewComponentBase("GitLab CLI", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type glabComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *glabComponent) GetAllVersions() ([]*gover.Version, error) {
	versionStrings, err := installer.Tools.GitLab.GetPackageReleases("gitlab-org/cli")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(versionStrings, versionRegexp, true)
}

func (c *glabComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("glab_%s_linux_%s.tar.gz", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, "v"+version.Raw, "downloads", fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "GitLab CLI"); err != nil {
		return err
	}
	// Extract the tar.gz file
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "/usr/local/bin/", true); err != nil {
		return err
	}

	// Set execute rights
	if err := os.Chmod("/usr/local/bin/glab", 0x755); err != nil {
		return err
	}
	return nil
}
