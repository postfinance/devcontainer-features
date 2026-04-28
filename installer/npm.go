package installer

import (
	"encoding/json"
	"fmt"

	"github.com/roemer/goext"
)

type npm struct{}

func (t npm) GetLatestPackageVersion(npmPackage string) (string, error) {
	stdout, stderr, err := goext.CmdRunners.Default.RunGetOutput("npm", "view", npmPackage, "version")
	if err != nil {
		return "", fmt.Errorf("failed getting latest version for %s: %s", npmPackage, stderr)
	}
	return stdout, nil
}

func (t npm) GetAllPackageVersions(npmPackage string) ([]string, error) {
	stdout, stderr, err := goext.CmdRunners.Default.RunGetOutput("npm", "view", npmPackage, "versions", "--json")
	if err != nil {
		return nil, fmt.Errorf("failed getting all versions for %s: %s", npmPackage, stderr)
	}
	var jsonData []string
	if err := json.Unmarshal([]byte(stdout), &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}

func (t npm) InstallGlobalPackageWithVersion(npmPackage, version string) error {
	packageWithVersion := npmPackage
	if version != "" && version != "latest" {
		packageWithVersion = fmt.Sprintf("%s@%s", npmPackage, version)
	}
	return t.InstallGlobalPackage(packageWithVersion)
}

func (t npm) InstallGlobalPackage(npmPackage string) error {
	_, stderr, err := goext.CmdRunners.Console.RunGetOutput("npm", "install", "-g", "--omit=dev", npmPackage)
	if err != nil {
		return fmt.Errorf("failed installing global npm package %s: %s", npmPackage, stderr)
	}
	return nil
}
