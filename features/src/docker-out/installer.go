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

var downloadRegistryBase = "https://download.docker.com"
var downloadRegistryPath = "/linux/static/stable/x86_64/"

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
	version := flag.String("version", "latest", "The version of the Docker CLI to install.")
	versionResolve := flag.Bool("versionResolve", false, "Whether to resolve the version to the latest available version.")
	composeVersion := flag.String("composeVersion", "latest", "The version of the Compose plugin to install.")
	composeVersionResolve := flag.Bool("composeVersionResolve", false, "Whether to resolve the version to the latest available version.")
	buildxVersion := flag.String("buildxVersion", "latest", "The version of the buildx plugin to install.")
	buildxVersionResolve := flag.Bool("buildxVersionResolve", false, "Whether to resolve the version to the latest available version.")

	flag.Parse()

	// Create and process the feature
	feature := installer.NewFeature("Docker-Out", false,
		&dockerCliComponent{
			ComponentBase: installer.NewComponentBase("Docker CLI", *version, *versionResolve),
		},
		&dockerComposeComponent{
			ComponentBase: installer.NewComponentBase("Docker Compose", *composeVersion, *composeVersionResolve),
		},
		&dockerBuildxComponent{
			ComponentBase: installer.NewComponentBase("Docker buildx", *buildxVersion, *buildxVersionResolve),
		},
	)
	return feature.Process()
}

//////////
// Implementation
//////////

type dockerCliComponent struct {
	*installer.ComponentBase
}

func (c *dockerCliComponent) GetAllVersions() ([]*gover.Version, error) {
	url, err := installer.Tools.Http.BuildUrl(downloadRegistryBase, downloadRegistryPath)
	if err != nil {
		return nil, err
	}
	versionRegexp := regexp.MustCompile(`(?m)^(\d+)\.(\d+)\.(\d+)$`)
	allVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndex(
		url,
		regexp.MustCompile(`^.*<a href="docker-([0-9\.]+).tgz">.*$`),
		versionRegexp)
	if err != nil {
		return nil, err
	}
	return allVersions, nil
}

func (c *dockerCliComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	fileName := fmt.Sprintf("docker-%s.tgz", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(downloadRegistryBase, downloadRegistryPath, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Docker CLI"); err != nil {
		return err
	}
	defer os.Remove(fileName)
	// Extract it
	tempDir, err := os.MkdirTemp("", "docker-extract")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)
	if err := installer.Tools.Compression.ExtractTarGz(fileName, tempDir, false); err != nil {
		return err
	}
	// Move the desired files
	if err := installer.Tools.FileSystem.MoveFile(filepath.Join(tempDir, "docker/docker"), "/usr/local/bin/docker"); err != nil {
		return err
	}

	// Copy the startup file
	if err := execr.Run(true, "cp", "-prf", "docker-init.sh", "/usr/local/share/docker-init.sh"); err != nil {
		return err
	}

	return nil
}

type dockerComposeComponent struct {
	*installer.ComponentBase
}

func (c *dockerComposeComponent) GetAllVersions() ([]*gover.Version, error) {
	versionRegexp := regexp.MustCompile(`(?m)^v(\d+)\.(\d+)\.(\d+)$`)
	versions := []*gover.Version{}
	allTags, err := installer.Tools.GitHub.GetTags("docker", "compose")
	if err != nil {
		return nil, err
	}
	for _, tag := range allTags {
		if versionRegexp.MatchString(tag.Name) {
			version, err := gover.ParseVersionFromRegex(tag.Name, versionRegexp)
			if err != nil {
				return nil, err
			}
			versions = append(versions, version)
		}
	}
	return versions, nil
}

func (c *dockerComposeComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	downloadUrl := fmt.Sprintf("https://github.com/docker/compose/releases/download/%s/docker-compose-linux-x86_64", version.Raw)
	if err := installer.Tools.Download.ToFile(downloadUrl, "/usr/local/lib/docker/cli-plugins/docker-compose", "Docker Compose"); err != nil {
		return err
	}
	// Apply executable permissions
	if err := execr.Run(true, "chmod", "+x", "/usr/local/lib/docker/cli-plugins/docker-compose"); err != nil {
		return err
	}
	return nil
}

type dockerBuildxComponent struct {
	*installer.ComponentBase
}

func (c *dockerBuildxComponent) GetAllVersions() ([]*gover.Version, error) {
	versionRegexp := regexp.MustCompile(`(?m)^v(\d+)\.(\d+)\.(\d+)$`)
	versions := []*gover.Version{}
	allTags, err := installer.Tools.GitHub.GetTags("docker", "buildx")
	if err != nil {
		return nil, err
	}
	for _, tag := range allTags {
		if versionRegexp.MatchString(tag.Name) {
			version, err := gover.ParseVersionFromRegex(tag.Name, versionRegexp)
			if err != nil {
				return nil, err
			}
			versions = append(versions, version)
		}
	}
	return versions, nil
}

func (c *dockerBuildxComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	downloadUrl := fmt.Sprintf("https://github.com/docker/buildx/releases/download/%s/buildx-%s.linux-amd64", version.Raw, version.Raw)
	if err := installer.Tools.Download.ToFile(downloadUrl, "/usr/local/lib/docker/cli-plugins/docker-buildx", "Docker buildx"); err != nil {
		return err
	}
	// Apply executable permissions
	if err := execr.Run(true, "chmod", "+x", "/usr/local/lib/docker/cli-plugins/docker-buildx"); err != nil {
		return err
	}
	return nil
}
