package main

import (
	"builder/installer"
	"fmt"
	"os"

	"github.com/roemer/gover"
)

//////////
// Main
//////////

func main() {
	if err := runMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	// Create and process the feature
	feature := installer.NewFeature("Make", true,
		&makeComponent{
			ComponentBase: installer.NewComponentBase("Make", installer.VERSION_SYSTEM_DEFAULT),
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type makeComponent struct {
	*installer.ComponentBase
}

func (c *makeComponent) InstallVersion(version *gover.Version) error {
	return installer.Tools.Apt.InstallDependencies("make")
}
