package docker

import (
	"builder/installer"
	"fmt"
	"regexp"

	"github.com/roemer/gover"
)

var dockerComposeVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

type DockerComposeComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *DockerComposeComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("docker", "compose")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, dockerComposeVersionRegexp, true)
}

func (c *DockerComposeComponent) InstallVersion(version *gover.Version) error {
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
	if err := installer.Tools.System.DownloadBinaryToDir(downloadUrl, "Docker Compose", "/usr/local/lib/docker/cli-plugins", "docker-compose"); err != nil {
		return err
	}
	return nil
}
