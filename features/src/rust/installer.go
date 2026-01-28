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

// Regex with 2-3 digits like 1.0 or 1.79.0
var threeDigitRegex *regexp.Regexp = regexp.MustCompile(`^?(\d+)\.(\d+)(?:\.(\d+))?$`)

// Full Regex versioning, like 1.0.0-alpha.2
var semVerRegex *regexp.Regexp = regexp.MustCompile(`^?(\d+)\.(\d+)(?:\.(\d+))?(?:-([a-z]+)(?:\.?(\d+))?)?$`)

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
	version := flag.String("version", "lts", "")
	rustupVersion := flag.String("rustupVersion", "latest", "")
	profile := flag.String("profile", "minimal", "")
	components := flag.String("components", "rustfmt,rust-analyzer,rust-src,clippy", "")
	enableWindowsTarget := flag.Bool("enableWindowsTarget", false, "")
	flag.Parse()

	// Create and process the feature
	feature := installer.NewFeature("Rust", true,
		&rustupComponent{
			ComponentBase: installer.NewComponentBase("rustup", *rustupVersion),
			profile:       *profile,
		},
		&rustComponent{
			ComponentBase: installer.NewComponentBase("rust", *version),
			components:    *components,
			profile:       *profile,
		},
		&buildEssentialComponent{
			ComponentBase: installer.NewComponentBase("build-essential", installer.VERSION_SYSTEM_DEFAULT),
		},
	)
	// Optional component
	if *enableWindowsTarget {
		feature.AddComponents(&windowsTargetComponent{
			ComponentBase: installer.NewComponentBase("windows-target", installer.VERSION_IRRELEVANT),
		})
	}
	// Last component
	feature.AddComponents(&permissionsComponent{
		ComponentBase: installer.NewComponentBase("permissions", installer.VERSION_IRRELEVANT),
	})
	return feature.Process()
}

//////////
// Implementation
//////////

type rustupComponent struct {
	*installer.ComponentBase
	profile string
}

func (c *rustupComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("rust-lang", "rustup")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, threeDigitRegex, true)
}

func (c *rustupComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}
	fileName := "rustup-init"
	downloadUrl := fmt.Sprintf("https://static.rust-lang.org/rustup/archive/%s/%s-unknown-linux-gnu/rustup-init", version.Raw, archPart)
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Rustup-Init"); err != nil {
		return err
	}
	// Install it
	if err := os.Chmod(fileName, os.ModePerm); err != nil {
		return err
	}
	if err := execr.Run(true, "./"+fileName, "-y", "--default-toolchain", "none", "--no-modify-path", "--profile", c.profile); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}

type rustComponent struct {
	*installer.ComponentBase
	profile    string
	components string
}

func (c *rustComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("rust-lang", "rust")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, semVerRegex, true)
}

func (c *rustComponent) InstallVersion(version *gover.Version) error {
	// Install it
	if err := execr.Run(true, "rustup", "toolchain", "install", "--profile", c.profile, "--no-self-update", version.Raw); err != nil {
		return err
	}
	// Installing the components
	fmt.Printf("Installing components: %s\n", c.components)
	args := []string{
		"component",
		"add",
	}
	for _, component := range strings.Split(c.components, ",") {
		trimmed := strings.TrimSpace(component)
		if trimmed != "" {
			args = append(args, trimmed)
		}
	}
	if err := execr.Run(true, "rustup", args...); err != nil {
		return err
	}
	return nil
}

type buildEssentialComponent struct {
	*installer.ComponentBase
}

func (c *buildEssentialComponent) InstallVersion(version *gover.Version) error {
	return installer.Tools.System.InstallPackages("build-essential")
}

type windowsTargetComponent struct {
	*installer.ComponentBase
}

func (c *windowsTargetComponent) InstallVersion(version *gover.Version) error {
	if err := execr.Run(true, "rustup", "target", "add", "x86_64-pc-windows-gnu"); err != nil {
		return err
	}
	if err := installer.Tools.System.InstallPackages("mingw-w64"); err != nil {
		return err
	}
	return nil
}

type permissionsComponent struct {
	*installer.ComponentBase
}

func (c *permissionsComponent) InstallVersion(version *gover.Version) error {
	return execr.Run(true, "chmod", "-R", "777", os.Getenv("RUSTUP_HOME"), os.Getenv("CARGO_HOME"))
}
