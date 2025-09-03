package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"

	"github.com/roemer/gover"
)

// ////////
// Main
// ////////

func main() {
	if err := runMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	// Create and process the feature
	yqVersion := flag.String("yqVersion", "system-default", "The version of yq to install.")
	gettextbaseVersion := flag.String("gettextbaseVersion", "system-default", "The version of gettext-base to install.")
	yamllintVersion := flag.String("yamllintVersion", "system-default", "The version of yamllint to install.")
	gitlfsVersion := flag.String("gitlfsVersion", "system-default", "The version of git-lfs to install.")
	sshpassVersion := flag.String("sshpassVersion", "system-default", "The version of sshpass to install.")

	feature := installer.NewFeature("CI Utility", true,
		&yqComponent{
			ComponentBase: installer.NewComponentBase("yq", *yqVersion),
		},
		&gettextbaseComponent{
			ComponentBase: installer.NewComponentBase("gettext-base", *gettextbaseVersion),
		},
		&yamllintComponent{
			ComponentBase: installer.NewComponentBase("yamllint", *yamllintVersion),
		},
		&gitlfsComponent{
			ComponentBase: installer.NewComponentBase("git-lfs", *gitlfsVersion),
		},
		&sshpassComponent{
			ComponentBase: installer.NewComponentBase("sshpass", *sshpassVersion),
		},
	)
	return feature.Process()
}

// ////////
// Implementation
// ////////

type yqComponent struct {
	*installer.ComponentBase
}

func (c *yqComponent) InstallVersion(*gover.Version) error {
	return installer.Tools.Apt.InstallDependencies("yq")
}

type gettextbaseComponent struct {
	*installer.ComponentBase
}

func (c *gettextbaseComponent) InstallVersion(*gover.Version) error {
	return installer.Tools.Apt.InstallDependencies("gettext-base")
}

type yamllintComponent struct {
	*installer.ComponentBase
}

func (c *yamllintComponent) InstallVersion(*gover.Version) error {
	return installer.Tools.Apt.InstallDependencies("yamllint")
}

type gitlfsComponent struct {
	*installer.ComponentBase
}

func (c *gitlfsComponent) InstallVersion(*gover.Version) error {
	return installer.Tools.Apt.InstallDependencies("git-lfs")
}

type sshpassComponent struct {
	*installer.ComponentBase
}

func (c *sshpassComponent) InstallVersion(*gover.Version) error {
	return installer.Tools.Apt.InstallDependencies("sshpass")
}
