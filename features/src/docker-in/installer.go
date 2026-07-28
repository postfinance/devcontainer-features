package main

import (
	"builder/installer"
	"builder/shared/docker"
	"bytes"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"text/template"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var dockerVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^(\d+)\.(\d+)\.(\d+)$`)

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

	installer.HandleOverride(downloadUrl, "https://download.docker.com/linux/static/stable", "docker-in-download-url")
	installer.HandleOverride(versionsUrl, "https://download.docker.com/linux/static/stable", "docker-in-versions-url")
	installer.HandleGitHubOverride(composeDownloadUrl, "docker/compose", "docker-in-compose-download-url")
	installer.HandleGitHubOverride(buildxDownloadUrl, "docker/buildx", "docker-in-buildx-download-url")
	installer.HandleOverride(configPath, "", "docker-in-config-path")

	// Create and process the feature
	feature := installer.NewFeature("Docker-In", false,
		&dockerComponent{
			ComponentBase: installer.NewComponentBase("Docker", *version),
			DownloadUrl:   *downloadUrl,
			VersionsUrl:   *versionsUrl,
			ConfigPath:    *configPath,
		},
		&docker.DockerComposeComponent{
			ComponentBase: installer.NewComponentBase("Docker Compose", *composeVersion),
			DownloadUrl:   *composeDownloadUrl,
		},
		&docker.DockerBuildxComponent{
			ComponentBase: installer.NewComponentBase("Docker buildx", *buildxVersion),
			DownloadUrl:   *buildxDownloadUrl,
		},
	)
	return feature.Process()
}

//////////
// Implementation
//////////

type dockerComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
	ConfigPath  string
}

func (c *dockerComponent) GetAllVersions() ([]*gover.Version, error) {
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
		dockerVersionRegexp)
	if err != nil {
		return nil, err
	}
	return allVersions, nil
}

func (c *dockerComponent) InstallVersion(version *gover.Version) error {
	// Install the system dependencies
	// XZ for alpine might be just xz
	if err := installer.Tools.System.InstallPackages("git", "procps", "iptables", "xz-utils"); err != nil {
		return err
	}
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
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Docker"); err != nil {
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
	// Install the binaries
	dockerFiles, err := os.ReadDir(filepath.Join(tempDir, "docker"))
	if err != nil {
		return err
	}
	for _, file := range dockerFiles {
		if file.IsDir() {
			continue
		}
		srcPath := filepath.Join(tempDir, "docker", file.Name())
		if err := installer.Tools.System.InstallBinaryToUsrLocalBin(srcPath, file.Name()); err != nil {
			return err
		}
	}
	// Generate the startup script
	t1, err := template.New("dind-init.sh.tmpl").ParseFiles("./dind-init.sh.tmpl")
	if err != nil {
		return err
	}
	data := struct {
		User string
	}{
		User: os.Getenv("_REMOTE_USER"),
	}
	var buf bytes.Buffer
	if err := t1.Execute(&buf, data); err != nil {
		return err
	}
	content := buf.String()
	if err := os.WriteFile("/usr/local/share/dind-init.sh", []byte(content), os.ModePerm); err != nil {
		return err
	}
	// Copy the default config.json
	if err := docker.InstallClientConfig(c.ConfigPath); err != nil {
		return err
	}

	return nil
}

func (c *dockerComponent) getArchitecturePathPart() (string, error) {
	return installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
}
