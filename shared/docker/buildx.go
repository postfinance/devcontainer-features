package docker

import (
	"builder/installer"
	"fmt"
	"regexp"

	"github.com/roemer/gover"
)

var dockerBuildxVersionRegexp *regexp.Regexp = regexp.MustCompile(`(?m)^v(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

type DockerBuildxComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *DockerBuildxComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("docker", "buildx")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, dockerBuildxVersionRegexp, true)
}

func (c *DockerBuildxComponent) InstallVersion(version *gover.Version) error {
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
	if err := installer.Tools.System.DownloadBinaryToDir(downloadUrl, "Docker buildx", "/usr/local/lib/docker/cli-plugins", "docker-buildx"); err != nil {
		return err
	}
	return nil
}
