package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var gitLfsVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

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

	installer.HandleGitHubOverride(downloadUrl, "git-lfs/git-lfs", "git-lfs-download-url")

	// Create and process the feature
	feature := installer.NewFeature("Git LFS", false,
		&gitLfsComponent{
			ComponentBase: installer.NewComponentBase("Git LFS", *version),
			DownloadUrl:   *downloadUrl,
		},
	)
	return feature.Process()
}

//////////
// Implementation
//////////

type gitLfsComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *gitLfsComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("git-lfs", "git-lfs")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, gitLfsVersionRegexp, true)
}

func (c *gitLfsComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	// https://github.com/git-lfs/git-lfs/releases/download/v3.7.0/git-lfs-linux-amd64-v3.7.0.tar.gz
	// https://github.com/git-lfs/git-lfs/releases/download/v3.7.0/git-lfs-linux-arm64-v3.7.0.tar.gz
	versionPart := fmt.Sprintf("v%s", version.Raw)
	fileName := fmt.Sprintf("git-lfs-linux-%s-%s.tar.gz", archPart, versionPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, versionPart, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Git LFS"); err != nil {
		return err
	}
	// Extract it
	tempDir, err := os.MkdirTemp("", "git-lfs-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, true); err != nil {
		return err
	}
	// Move the desired files
	if err := installer.Tools.FileSystem.MoveFile(filepath.Join(tempDir, "git-lfs"), "/usr/local/bin/git-lfs"); err != nil {
		return err
	}
	// Apply executable permissions
	if err := execr.Run(true, "chmod", "+x", "/usr/local/bin/git-lfs"); err != nil {
		return err
	}
	// Install
	if err := execr.Run(true, "git", "lfs", "install"); err != nil {
		return err
	}
	return nil
}
