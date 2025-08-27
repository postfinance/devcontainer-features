package installer

import (
	"encoding/json"
	"fmt"

	"github.com/roemer/gotaskr/execr"
)

type npm struct{}

func (npm) GetLatestPackageVersion(npmPackage string) (string, error) {
	stdout, stderr, err := execr.RunGetOutput(false, "npm", "view", npmPackage, "version")
	if err != nil {
		return "", fmt.Errorf("failed getting latest version for %s: %s", npmPackage, stderr)
	}
	return stdout, nil
}

func (npm) GetAllPackageVersions(npmPackage string) ([]string, error) {
	stdout, stderr, err := execr.RunGetOutput(false, "npm", "view", npmPackage, "versions", "--json")
	if err != nil {
		return nil, fmt.Errorf("failed getting all versions for %s: %s", npmPackage, stderr)
	}
	var jsonData []string
	if err := json.Unmarshal([]byte(stdout), &jsonData); err != nil {
		return nil, err
	}
	return jsonData, nil
}
