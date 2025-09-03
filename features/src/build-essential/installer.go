package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

const packageName = "build-essential"

var versionOnlyRegex *regexp.Regexp = regexp.MustCompile(`[^\|]*\|\s+([^\|]*)\s+\|.*$`)
var buildEssentialVersionRegex *regexp.Regexp = regexp.MustCompile(`^(?P<d1>\d+)\.(?P<d2>\d+)`)

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
	// Handle the flags
	version := flag.String("version", "latest", "The version of build-essential to install.")
	flag.Parse()

	// Create and process the feature
	feature := installer.NewFeature("Build Essential", true,
		&buildEssentialComponent{
			ComponentBase: installer.NewComponentBase(packageName, *version),
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type buildEssentialComponent struct {
	*installer.ComponentBase
}

func (c *buildEssentialComponent) GetAllVersions() ([]*gover.Version, error) {
	if err := execr.Run(true, "apt-get", "update"); err != nil {
		return []*gover.Version{}, err
	}

	str1, _, err := execr.RunGetOutput(true, "apt-cache", "madison", packageName)
	if err != nil {
		return []*gover.Version{}, err
	}

	allFiles := strings.Split(str1, "\n")
	allVersions := []*gover.Version{}
	for _, item := range allFiles {
		matches := versionOnlyRegex.FindStringSubmatch(item)
		if matches == nil {
			continue
		}
		allVersions = append(allVersions, gover.MustParseVersionFromRegex(matches[1], buildEssentialVersionRegex))
	}

	return allVersions, nil
}

func (c *buildEssentialComponent) InstallVersion(version *gover.Version) error {
	if version == gover.EmptyVersion {
		return installer.Tools.Apt.InstallDependencies(packageName, "pkg-config", "libpcsclite-dev")
	} else {
		return installer.Tools.Apt.InstallDependencies(fmt.Sprintf(packageName+"=%s", version.Raw), "pkg-config", "libpcsclite-dev")
	}
}
