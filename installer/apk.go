package installer

import (
	"strings"

	"github.com/roemer/gotaskr/execr"
)

type apk struct{}

func (a apk) InstallDependencies(dependencies ...string) error {
	args := append([]string{"add", "--no-cache"}, dependencies...)
	if err := execr.Run(true, "apk", args...); err != nil {
		return err
	}
	return nil
}

func (a apk) InstallLocalPackage(packagePath string) error {
	if !strings.HasPrefix(packagePath, "./") {
		packagePath = "./" + packagePath
	}
	if err := execr.Run(true, "apk", "add", "--no-cache", packagePath); err != nil {
		return err
	}
	return nil
}
