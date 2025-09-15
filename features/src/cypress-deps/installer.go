package main

import (
	"builder/installer"
	"fmt"
	"os"
)

func main() {
	if err := runMain(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func runMain() error {
	fmt.Println("Installing Cypress Dependencies")
	osInfo, err := installer.Tools.System.GetOsInfo()
	if err != nil {
		return fmt.Errorf("failed to get OS info: %w", err)
	}

	// Common dependencies for all supported distros
	commonDeps := []string{
		"libgbm-dev",
		"libnotify-dev",
		"libnss3",
		"libxss1",
		"libxtst6",
		"xauth",
		"xvfb",
	}

	var deps []string

	debianTrixieOrNewer := osInfo.Vendor == "debian" && func() bool {
		var major int
		fmt.Sscanf(osInfo.VersionId, "%d", &major)
		return major >= 13
	}()

	ubuntuNobleOrNewer := osInfo.Vendor == "ubuntu" && func() bool {
		var major int
		fmt.Sscanf(osInfo.VersionId, "%d", &major)
		return major >= 24
	}()

	if debianTrixieOrNewer || ubuntuNobleOrNewer {
		deps = append([]string{
			"libgtk2.0-0t64",
			"libgtk-3-0t64",
			"libasound2t64",
		}, commonDeps...)
	} else {
		deps = append([]string{
			"libgtk2.0-0",
			"libgtk-3-0",
			"libasound2",
		}, commonDeps...)
	}

	return installer.Tools.Apt.InstallDependencies(deps...)
}
