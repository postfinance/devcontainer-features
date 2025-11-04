package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var dockerCliVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^(\d+)\.(\d+)\.(\d+)$`)
var dockerComposeVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(?P<raw>(\d+)\.(\d+)\.(\d+))$`)
var dockerBuildxVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

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
	composeVersion := flag.String("composeVersion", "latest", "")
	buildxVersion := flag.String("buildxVersion", "latest", "")
	buildxDownloadUrl := flag.String("buildxDownloadUrl", "", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	versionsUrl := flag.String("versionsUrl", "", "")
	composeDownloadUrl := flag.String("composeDownloadUrl", "", "")
	configPath := flag.String("configPath", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrl, "https://download.docker.com/linux/static/stable", "docker-out-download-url")
	installer.HandleOverride(versionsUrl, "https://download.docker.com/linux/static/stable", "docker-out-versions-url")
	installer.HandleGitHubOverride(composeDownloadUrl, "docker/compose", "docker-out-compose-download-url")
	installer.HandleGitHubOverride(buildxDownloadUrl, "docker/buildx", "docker-out-buildx-download-url")
	installer.HandleOverride(configPath, "", "docker-out-config-path")

	// Create and process the feature
	feature := installer.NewFeature("Docker-Out", false,
		&dockerCliComponent{
			ComponentBase: installer.NewComponentBase("Docker CLI", *version),
			DownloadUrl:   *downloadUrl,
			VersionsUrl:   *versionsUrl,
			ConfigPath:    *configPath,
		},
		&dockerComposeComponent{
			ComponentBase: installer.NewComponentBase("Docker Compose", *composeVersion),
			DownloadUrl:   *composeDownloadUrl,
		},
		&dockerBuildxComponent{
			ComponentBase: installer.NewComponentBase("Docker buildx", *buildxVersion),
			DownloadUrl:   *buildxDownloadUrl,
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
	DownloadUrl string
	VersionsUrl string
	ConfigPath  string
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
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, architecturePathPart, fileName)
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
	// Copy the default config.json
	if c.ConfigPath != "" {
		fileContent, err := installer.ReadConfigFile(c.ConfigPath)
		if err != nil {
			return err
		}
		userName := os.Getenv("_REMOTE_USER")
		dockerDir := fmt.Sprintf("/home/%s/.docker", userName)
		configDest := fmt.Sprintf("%s/config.json", dockerDir)
		// Ensure directory exists
		if err := os.MkdirAll(dockerDir, 0700); err != nil {
			return err
		}
		// Write config file
		if err := os.WriteFile(configDest, fileContent, 0600); err != nil {
			return err
		}
		// Set ownership
		usr, err := user.Lookup(userName)
		if err != nil {
			return err
		}
		uid, _ := strconv.Atoi(usr.Uid)
		gid, _ := strconv.Atoi(usr.Gid)
		if err := os.Chown(dockerDir, uid, gid); err != nil {
			return err
		}
		if err := os.Chown(configDest, uid, gid); err != nil {
			return err
		}
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
	DownloadUrl string
}

func (c *dockerComposeComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("docker", "compose")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, dockerComposeVersionRegexp, true)
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
	// https://github.com/docker/compose/releases/download/v2.39.2/docker-compose-linux-x86_64
	versionPart := fmt.Sprintf("v%s", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, versionPart, fmt.Sprintf("docker-compose-linux-%s", archPart))
	if err != nil {
		return err
	}
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
	DownloadUrl string
}

func (c *dockerBuildxComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("docker", "buildx")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, dockerBuildxVersionRegexp, true)
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
	// https://github.com/docker/buildx/releases/download/v0.26.1/buildx-v0.26.1.linux-amd64
	versionPart := fmt.Sprintf("v%s", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, versionPart, fmt.Sprintf("buildx-%s.linux-%s", versionPart, archPart))
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, "/usr/local/lib/docker/cli-plugins/docker-buildx", "Docker buildx"); err != nil {
		return err
	}
	// Apply executable permissions
	if err := execr.Run(true, "chmod", "+x", "/usr/local/lib/docker/cli-plugins/docker-buildx"); err != nil {
		return err
	}
	return nil
}
