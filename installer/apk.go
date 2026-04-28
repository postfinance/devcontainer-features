package installer

import (
	"fmt"
	"strings"

	"github.com/roemer/goext"
)

type apk struct{}

func (a apk) InstallDependencies(dependencies ...string) error {
	args := append([]string{"add", "--no-cache"}, dependencies...)
	if err := goext.CmdRunners.Console.Run("apk", args...); err != nil {
		return fmt.Errorf("failed to install dependencies: %w", err)
	}
	return nil
}

func (a apk) InstallLocalPackage(packagePath string) error {
	if !strings.HasPrefix(packagePath, "./") {
		packagePath = "./" + packagePath
	}
	if err := goext.CmdRunners.Console.Run("apk", "add", "--no-cache", packagePath); err != nil {
		return fmt.Errorf("failed to install local package: %w", err)
	}
	return nil
}
