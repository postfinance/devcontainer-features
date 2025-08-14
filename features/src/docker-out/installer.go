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

var dockerCliVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^(\d+)\.(\d+)\.(\d+)$`)
var dockerComposeVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(\d+)\.(\d+)\.(\d+)$`)
var dockerBuildxVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(\d+)\.(\d+)\.(\d+)$`)

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
	composeVersion := flag.String("composeVersion", "latest", "")
	composeVersionResolve := flag.Bool("composeVersionResolve", false, "")
	buildxVersion := flag.String("buildxVersion", "latest", "")
	buildxVersionResolve := flag.Bool("buildxVersionResolve", false, "")
	downloadUrlBase := flag.String("downloadUrlBase", "", "")
	downloadUrlPath := flag.String("downloadUrlPath", "", "")
	versionsUrl := flag.String("versionsUrl", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrlBase, "https://download.docker.com", "docker-out-download-url-base")
	installer.HandleOverride(downloadUrlPath, "/linux/static/stable", "docker-out-download-url-path")
	installer.HandleOverride(versionsUrl, "https://download.docker.com/linux/static/stable", "docker-out-versions-url")

	// Create and process the feature
	feature := installer.NewFeature("Docker-Out", false,
		&dockerCliComponent{
			ComponentBase:   installer.NewComponentBase("Docker CLI", *version, *versionResolve),
			DownloadUrlBase: *downloadUrlBase,
			DownloadUrlPath: *downloadUrlPath,
			VersionsUrl:     *versionsUrl,
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

// Docker CLI

type dockerCliComponent struct {
	*installer.ComponentBase
	DownloadUrlBase string
	DownloadUrlPath string
	VersionsUrl     string
}

func (c *dockerCliComponent) GetAllVersions() ([]*gover.Version, error) {
	// Download the file
	architecturePathPart, err := c.getArchitecturePathPart()
	if err != nil {
		return nil, err
	}
	url, err := installer.Tools.Http.BuildUrl(c.VersionsUrl, architecturePathPart)
	if err != nil {
		return nil, err
	}
	allVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndex(
		url,
		regexp.MustCompile(`^.*<a href="docker-([0-9\.]+).tgz">.*$`),
		dockerCliVersionRegexp)
	if err != nil {
		return nil, err
	}
	return allVersions, nil
}

func (c *dockerCliComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	architecturePathPart, err := c.getArchitecturePathPart()
	if err != nil {
		return err
	}
	fileName := fmt.Sprintf("docker-%s.tgz", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrlBase, c.DownloadUrlPath, architecturePathPart, fileName)
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

func (c *dockerCliComponent) getArchitecturePathPart() (string, error) {
	return installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
}

// Docker Compose

type dockerComposeComponent struct {
	*installer.ComponentBase
}

func (c *dockerComposeComponent) GetAllVersions() ([]*gover.Version, error) {
	versions := []*gover.Version{}
	allTags, err := installer.Tools.GitHub.GetTags("docker", "compose")
	if err != nil {
		return nil, err
	}
	for _, tag := range allTags {
		if dockerComposeVersionRegexp.MatchString(tag.Name) {
			version, err := gover.ParseVersionFromRegex(tag.Name, dockerComposeVersionRegexp)
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
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}
	downloadUrl := fmt.Sprintf("https://github.com/docker/compose/releases/download/%s/docker-compose-linux-%s", version.Raw, archPart)
	if err := installer.Tools.Download.ToFile(downloadUrl, "/usr/local/lib/docker/cli-plugins/docker-compose", "Docker Compose"); err != nil {
		return err
	}
	// Apply executable permissions
	if err := execr.Run(true, "chmod", "+x", "/usr/local/lib/docker/cli-plugins/docker-compose"); err != nil {
		return err
	}
	return nil
}

// Docker buildx

type dockerBuildxComponent struct {
	*installer.ComponentBase
}

func (c *dockerBuildxComponent) GetAllVersions() ([]*gover.Version, error) {
	versions := []*gover.Version{}
	allTags, err := installer.Tools.GitHub.GetTags("docker", "buildx")
	if err != nil {
		return nil, err
	}
	for _, tag := range allTags {
		if dockerBuildxVersionRegexp.MatchString(tag.Name) {
			version, err := gover.ParseVersionFromRegex(tag.Name, dockerBuildxVersionRegexp)
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
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	downloadUrl := fmt.Sprintf("https://github.com/docker/buildx/releases/download/%s/buildx-%s.linux-%s", version.Raw, version.Raw, archPart)
	if err := installer.Tools.Download.ToFile(downloadUrl, "/usr/local/lib/docker/cli-plugins/docker-buildx", "Docker buildx"); err != nil {
		return err
	}
	// Apply executable permissions
	if err := execr.Run(true, "chmod", "+x", "/usr/local/lib/docker/cli-plugins/docker-buildx"); err != nil {
		return err
	}
	return nil
}
