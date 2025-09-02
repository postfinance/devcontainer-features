package installer

import (
	"fmt"
	"os"
	"strings"
)

type apt struct{}

var Apt apt = apt{}

func (a apt) InstallDependencies(dependencies ...string) error {
	if err := Run(false, "apt-get", "update"); err != nil {
		return err
	}

	args := append([]string{"install", "-y"}, dependencies...)
	if err := Run(true, "apt-get", args...); err != nil {
		return err
	}

	a.CleanCache()
	return nil
}

func (a apt) InstallLocalPackage(packagePath string) error {
	if !strings.HasPrefix(packagePath, "./") {
		packagePath = "./" + packagePath
	}
	if err := Run(false, "apt-get", "update"); err != nil {
		return err
	}
	if err := Run(true, "apt-get", "install", "-y", packagePath); err != nil {
		return err
	}
	a.CleanCache()
	return nil
}

func (apt) CleanCache() {
	fmt.Println("Cleaning apt cache")
	os.RemoveAll("/var/lib/apt/lists")
}
