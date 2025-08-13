package main

import (
	"builder/installer"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"runtime"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

const versionsIndexUrl = "https://ziglang.org/download/index.json"

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
	version := flag.String("version", "latest", "The version of Zig to install.")
	isExactVersion := flag.Bool("isExactVersion", false, "Whether to install the exact version specified.")
	downloadRegistryBase := flag.String("downloadRegistryBase", "", "The download registry to use for Zig binaries.")
	downloadRegistryPath := flag.String("downloadRegistryPath", "", "The download registry path to use for Zig binaries.")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadRegistryBase, "https://ziglang.org", "zig-download-registry-base")
	installer.HandleOverride(downloadRegistryPath, "/download", "zig-download-registry-path")

	// Create and process the feature
	feature := installer.NewFeature("Zig", true,
		&zigComponent{
			ComponentBase:        installer.NewComponentBase("Zig", *version, *isExactVersion),
			DownloadRegistryBase: *downloadRegistryBase,
			DownloadRegistryPath: *downloadRegistryPath,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type zigComponent struct {
	*installer.ComponentBase
	DownloadRegistryBase string
	DownloadRegistryPath string
}

func (c *zigComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(versionsIndexUrl)
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
	var fileName string
	// The url format changed from OS-Arch-Version to Arch-OS-Version with 0.14.1
	if version.LessThan(gover.ParseSimple(0, 14, 1)) {
		switch runtime.GOARCH {
		case "amd64":
			fileName = fmt.Sprintf("zig-linux-x86_64-%s.tar.xz", version.Raw)
		case "arm64":
			fileName = fmt.Sprintf("zig-linux-aarch64-%s.tar.xz", version.Raw)
		default:
			return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
		}
	} else {
		switch runtime.GOARCH {
		case "amd64":
			fileName = fmt.Sprintf("zig-x86_64-linux-%s.tar.xz", version.Raw)
		case "arm64":
			fileName = fmt.Sprintf("zig-aarch64-linux-%s.tar.xz", version.Raw)
		default:
			return fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
		}
	}

	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadRegistryBase, c.DownloadRegistryPath, version.Raw, fileName)
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
