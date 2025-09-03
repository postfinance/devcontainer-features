package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/roemer/gover"
)

func main() {
	if err := runMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	packagesFlag := flag.String("packages", "", "Comma-separated list of system packages to install.")
	flag.Parse()

	packages := parsePackages(*packagesFlag)
	if len(packages) == 0 {
		fmt.Println("No packages specified for installation.")
		return nil
	}

	feature := installer.NewFeature("System Packages", true,
		&systemPackagesComponent{
			ComponentBase: installer.NewComponentBase("System Packages", "system-default"),
			Packages:      packages,
		})
	return feature.Process()
}

type systemPackagesComponent struct {
	*installer.ComponentBase
	Packages []string
}

func (c *systemPackagesComponent) InstallVersion(version *gover.Version) error {
	return installer.Tools.System.InstallPackages(c.Packages)
}

func parsePackages(flagValue string) []string {
	if flagValue == "" {
		return nil
	}
	parts := strings.Split(flagValue, ",")
	var pkgs []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			pkgs = append(pkgs, trimmed)
		}
	}
	return pkgs
}
