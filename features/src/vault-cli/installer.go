package main

import (
	"builder/installer"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var versionRegexp *regexp.Regexp = regexp.MustCompile(`(?m:)^(\d+)\.(\d+)\.(\d+)$`)

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

	installer.HandleOverride(downloadUrlBase, "https://releases.hashicorp.com", "vault-cli-download-url-base")
	installer.HandleOverride(downloadUrlPath, "/vault", "vault-cli-download-url-path")
	installer.HandleOverride(versionsUrl, "https://releases.hashicorp.com/vault/index.json", "vault-cli-versions-url")

	// Create and process the feature
	feature := installer.NewFeature("Vault CLI", true,
		&vaultCliComponent{
			ComponentBase:   installer.NewComponentBase("Vault CLI", *version, *versionResolve),
			DownloadUrlBase: *downloadUrlBase,
			DownloadUrlPath: *downloadUrlPath,
			VersionsUrl:     *versionsUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type vaultCliComponent struct {
	*installer.ComponentBase
	DownloadUrlBase string
	DownloadUrlPath string
	VersionsUrl     string
}

func (c *vaultCliComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(c.VersionsUrl)
	if err != nil {
		return nil, err
	}
	var jsonData map[string]any
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	versionsObj, ok := jsonData["versions"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("versions not found in vault index json")
	}
	versions := []*gover.Version{}
	for versionString := range versionsObj {
		version, err := gover.ParseVersionFromRegex(versionString, versionRegexp)
		if err != nil {
			continue // skip invalid
		}
		versions = append(versions, version)
	}
	return versions, nil
}

func (c *vaultCliComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("vault_%s_linux_%s.zip", version.Raw, archPart)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrlBase, c.DownloadUrlPath, version.Raw, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Vault CLI"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractZip(fileName, "vault", false); err != nil {
		return err
	}
	// Install
	configureCmd := exec.Command("install", "-m", "0755", "vault/vault", "/usr/local/bin/vault")
	if err := execr.RunCommand(true, configureCmd); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll("vault"); err != nil {
		return err
	}
	return nil
}
