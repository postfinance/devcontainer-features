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

var versionRegex *regexp.Regexp = regexp.MustCompile(`(?m:)^(\d+).(\d+)(?:.(\d+)(?:.(\d+))?)?(?:-(rc\d+))?$`)

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
	includeJre := flag.Bool("includeJre", true, "")
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file (global/per-feature overrides)
	if err := installer.LoadOverrides(); err != nil {
		return err
	}
	// Handle overrides and defaults
	installer.HandleOverride(downloadUrl, "https://binaries.sonarsource.com/Distribution/sonar-scanner-cli", "sonar-scanner-cli-download-url")

	// Create and process the feature
	feature := installer.NewFeature("SonarScanner CLI", true,
		&sonarscannerCliComponent{
			ComponentBase: installer.NewComponentBase("SonarScanner CLI", *version),
			IncludeJre:    *includeJre,
			DownloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type sonarscannerCliComponent struct {
	*installer.ComponentBase
	IncludeJre  bool
	DownloadUrl string
}

func (c *sonarscannerCliComponent) IsFullVersion(referenceVersion *gover.Version) bool {
	return len(referenceVersion.Segments) == 4 && referenceVersion.DefinedSegmentCount() == 4
}

func (c *sonarscannerCliComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("SonarSource", "sonar-scanner-cli")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, versionRegex, true)
}

func (c *sonarscannerCliComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}
	var fileName string
	if c.IncludeJre {
		fileName = fmt.Sprintf("sonar-scanner-cli-%s-linux-%s.zip", version.Raw, archPart)
	} else {
		fileName = fmt.Sprintf("sonar-scanner-cli-%s.zip", version.Raw)
	}
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Sonar"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractZip(fileName, "/usr/local/sonar-scanner", true); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}
