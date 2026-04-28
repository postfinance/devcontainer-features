package installer

import (
	"fmt"
	"os"
	"strings"

	"github.com/roemer/goext"
)

type apt struct{}

func (a apt) Update() error {
	return goext.CmdRunners.Default.Run("apt-get", "update")
}

func (a apt) InstallDependencies(dependencies ...string) error {
	if len(dependencies) == 0 {
		return nil
	}
	if err := a.Update(); err != nil {
		return fmt.Errorf("failed to update apt-get: %w", err)
	}
	args := append([]string{"install", "-y"}, dependencies...)
	if err := goext.CmdRunners.Console.Run("apt-get", args...); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	a.CleanCache()
	return nil
}

func (a apt) InstallLocalPackage(packagePath string) error {
	if !strings.HasPrefix(packagePath, "./") {
		packagePath = "./" + packagePath
	}
	if err := a.Update(); err != nil {
		return fmt.Errorf("failed to update apt-get: %w", err)
	}
	if err := goext.CmdRunners.Console.Run("apt-get", "install", "-y", packagePath); err != nil {
		return fmt.Errorf("failed to install local package: %w", err)
	}
	a.CleanCache()
	return nil
}

func (apt) CleanCache() {
	fmt.Println("Cleaning apt cache")
	os.RemoveAll("/var/lib/apt/lists")
}
