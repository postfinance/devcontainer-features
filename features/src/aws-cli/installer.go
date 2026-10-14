package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/roemer/goext"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var versionRegex *regexp.Regexp = regexp.MustCompile(`(?m)^(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

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
	installer.HandleOverride(downloadUrl, "https://awscli.amazonaws.com", "aws-cli-download-url")

	// Create and process the feature
	feature := installer.NewFeature("AWS CLI", false,
		&awsCliComponent{
			ComponentBase: installer.NewComponentBase("AWS CLI", *version),
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type awsCliComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *awsCliComponent) GetAllVersions() ([]*gover.Version, error) {
	tags, err := installer.Tools.GitHub.GetTags("aws", "aws-cli")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(tags, versionRegex, true)
}

func (c *awsCliComponent) InstallVersion(version *gover.Version) error {
	// Map architecture
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}

	// Build the download URL and file name
	// https://awscli.amazonaws.com/awscli-exe-linux-x86_64-2.22.35.zip
	fileName := fmt.Sprintf("awscli-exe-linux-%s-%s.zip", archPart, version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, fileName)
	if err != nil {
		return err
	}

	// Download the file
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "AWS CLI"); err != nil {
		return err
	}
	defer os.Remove(fileName)

	// Extract to temp directory
	tempDir, err := os.MkdirTemp("", "aws-cli-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	if err := installer.Tools.Compression.ExtractZip(fileName, tempDir, false); err != nil {
		return err
	}

	// Run the AWS installer script
	awsInstallerPath := fmt.Sprintf("%s/aws/install", tempDir)
	if err := goext.CmdRunners.Console.Run(awsInstallerPath, "--bin-dir", "/usr/local/bin", "--install-dir", "/usr/local/aws-cli", "--update"); err != nil {
		return err
	}

	return nil
}
