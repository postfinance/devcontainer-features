package main

import (
	"builder/installer"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var versionRegexp *regexp.Regexp = gover.RegexpSimple

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
	versionResolve := flag.Bool("versionResolve", false, "")
	downloadUrlBase := flag.String("downloadUrlBase", "", "")
	downloadUrlPath := flag.String("downloadUrlPath", "", "")
	versionsUrl := flag.String("versionsUrl", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrlBase, "https://ziglang.org", "zig-download-url-base")
	installer.HandleOverride(downloadUrlPath, "/download", "zig-download-url-path")
	installer.HandleOverride(versionsUrl, "https://ziglang.org/download/index.json", "zig-versions-url")

	// Create and process the feature
	feature := installer.NewFeature("Zig", true,
		&zigComponent{
			ComponentBase:   installer.NewComponentBase("Zig", *version, *versionResolve),
			DownloadUrlBase: *downloadUrlBase,
			DownloadUrlPath: *downloadUrlPath,
			VersionsUrl:     *versionsUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type zigComponent struct {
	*installer.ComponentBase
	DownloadUrlBase string
	DownloadUrlPath string
	VersionsUrl     string
}

func (c *zigComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(c.VersionsUrl)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]any
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	versions := []*gover.Version{}
	for key := range jsonData {
		if key == "master" {
			continue
		}
		versionString := key
		version, err := gover.ParseVersionFromRegex(versionString, versionRegexp)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}

	return versions, nil
}

func (c *zigComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}
	var fileName string
	// The url format changed from OS-Arch-Version to Arch-OS-Version with 0.14.1
	if version.LessThan(gover.ParseSimple(0, 14, 1)) {
		fileName = fmt.Sprintf("zig-linux-%s-%s.tar.xz", archPart, version.Raw)
	} else {
		fileName = fmt.Sprintf("zig-%s-linux-%s.tar.xz", archPart, version.Raw)
	}
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrlBase, c.DownloadUrlPath, version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "zig"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractTarXz(fileName, "/usr/local/zig", true); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}
