package installer

import (
	"fmt"
	"os"
	"strings"

	"github.com/roemer/gotaskr/execr"
)

type apt struct{}

func (a apt) InstallDependencies(dependencies ...string) error {
	if len(dependencies) == 0 {
		return nil
	}
	if err := execr.Run(false, "apt-get", "update"); err != nil {
		return fmt.Errorf("failed to update apt-get: %w", err)
	}
	args := append([]string{"install", "-y"}, dependencies...)
	if err := execr.Run(true, "apt-get", args...); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	a.CleanCache()
	return nil
}

func (a apt) InstallLocalPackage(packagePath string) error {
	if !strings.HasPrefix(packagePath, "./") {
		packagePath = "./" + packagePath
	}
	if err := execr.Run(false, "apt-get", "update"); err != nil {
		return fmt.Errorf("failed to update apt-get: %w", err)
	}
	if err := execr.Run(true, "apt-get", "install", "-y", packagePath); err != nil {
		return fmt.Errorf("failed to install local package: %w", err)
	}
	a.CleanCache()
	return nil
}

func (apt) CleanCache() {
	fmt.Println("Cleaning apt cache")
	os.RemoveAll("/var/lib/apt/lists")
}
