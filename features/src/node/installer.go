package main

import (
	"builder/installer"
	"encoding/json"
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

var nodeVersionRegex *regexp.Regexp = regexp.MustCompile(`(?m:^v(?P<raw>(\d+)(?:\.(\d+))?(?:\.(\d+)))?$)`)
var npmVersionRegex *regexp.Regexp = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)(?:-(.+?))?(?:[-\.](\d+))?(?:[-\.](\d+))?$`)

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
	version := flag.String("version", "lts", "")
	npmVersion := flag.String("npmVersion", "included", "")
	yarnVersion := flag.String("yarnVersion", "none", "")
	pnpmVersion := flag.String("pnpmVersion", "none", "")
	corepackVersion := flag.String("corepackVersion", "none", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	versionsUrl := flag.String("versionsUrl", "", "")
	globalNpmRegistry := flag.String("globalNpmRegistry", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrl, "https://nodejs.org/dist", "node-download-url")
	installer.HandleOverride(versionsUrl, "https://nodejs.org/dist/index.json", "node-versions-url")

	// Create and process the feature
	feature := installer.NewFeature("Node.js", false,
		&nodeComponent{
			ComponentBase:     installer.NewComponentBase("Node.js", *version),
			DownloadUrl:       *downloadUrl,
			VersionsUrl:       *versionsUrl,
			GlobalNpmRegistry: *globalNpmRegistry,
		},
		&npmComponent{
			ComponentBase: installer.NewComponentBase("NPM", *npmVersion),
			PackageName:   "npm",
		},
		&npmComponent{
			ComponentBase: installer.NewComponentBase("Yarn", *yarnVersion),
			PackageName:   "yarn",
		},
		&npmComponent{
			ComponentBase: installer.NewComponentBase("Pnpm", *pnpmVersion),
			PackageName:   "pnpm",
		},
		&npmComponent{
			ComponentBase: installer.NewComponentBase("corepack", *corepackVersion),
			PackageName:   "corepack",
		},
	)
	return feature.Process()
}

//////////
// Implementation
//////////

type nodeComponent struct {
	*installer.ComponentBase
	DownloadUrl       string
	VersionsUrl       string
	GlobalNpmRegistry string
}

func (c *nodeComponent) GetAllVersions() ([]*gover.Version, error) {
	ltsVersionsOnly := c.GetRequestedVersion() == installer.VERSION_LTS
	versionFileContent, err := installer.Tools.Download.AsBytes(c.VersionsUrl)
	if err != nil {
		return nil, err
	}
	var jsonData []*struct {
		Version string      `json:"version"`
		Lts     interface{} `json:"lts"` // Can be a boolean or a string
	}
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	for _, entry := range jsonData {
		version := gover.MustParseVersionFromRegex(entry.Version, nodeVersionRegex)
		isLtsEntry := false
		switch entry.Lts.(type) {
		case bool:
			isLtsEntry = entry.Lts.(bool)
		case string:
			isLtsEntry = true
		}
		if ltsVersionsOnly && !isLtsEntry {
			continue
		}
		allVersions = append(allVersions, version)
	}
	return allVersions, nil
}

func (c *nodeComponent) InstallVersion(version *gover.Version) error {
	// Prepare and download node
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	versionString := "v" + version.Raw
	releaseName := fmt.Sprintf("node-%s-linux-%s", versionString, archPart)
	fileName := fmt.Sprintf("%s.tar.gz", releaseName)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, versionString, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Node.JS"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "./", false); err != nil {
		return err
	}
	// Move the needed folders
	if err := installer.Tools.FileSystem.MoveFolders([]string{
		filepath.Join(releaseName, "bin"),
		filepath.Join(releaseName, "include"),
		filepath.Join(releaseName, "lib"),
		filepath.Join(releaseName, "share"),
	}, "/usr"); err != nil {
		return err
	}
	// Configure NPM
	if c.GlobalNpmRegistry != "" {
		if err := execr.Run(true, "npm", "set", "--global", "registry", c.GlobalNpmRegistry); err != nil {
			return err
		}
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll(releaseName); err != nil {
		return err
	}
	return nil
}

type npmComponent struct {
	*installer.ComponentBase
	PackageName string
}

func (c *npmComponent) GetAllVersions() ([]*gover.Version, error) {
	allVersionStrings, err := installer.Tools.Npm.GetAllPackageVersions(c.PackageName)
	if err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	for _, versionString := range allVersionStrings {
		version, err := gover.ParseVersionFromRegex(versionString, npmVersionRegex)
		if err != nil {
			return nil, err
		}
		allVersions = append(allVersions, version)
	}
	return allVersions, nil
}

func (c *npmComponent) GetLatestVersion() (*gover.Version, error) {
	latestVersion, err := installer.Tools.Npm.GetLatestPackageVersion(c.PackageName)
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, npmVersionRegex)
	if err != nil {
		return nil, err
	}
	return version, err
}

func (c *npmComponent) InstallVersion(version *gover.Version) error {
	return execr.Run(true, "npm", "install", "--global", c.PackageName+"@"+version.Raw)
}
