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
	fmt.Println("Installing Eclipse Dependencies")

	osInfo := installer.Tools.System.OsInfo()

	debianTrixieOrNewer := osInfo.IsDebian() && osInfo.MajorVersion() >= 13
	ubuntuNobleOrNewer := osInfo.IsUbuntu() && osInfo.MajorVersion() >= 24

	webkitPackage := "libwebkit2gtk-4.0-37"
	if debianTrixieOrNewer || ubuntuNobleOrNewer {
		webkitPackage = "libwebkit2gtk-4.1-0"
	}

	return installer.Tools.Apt.InstallDependencies(
		"libswt-gtk-4-jni",
		webkitPackage,
	)
}
