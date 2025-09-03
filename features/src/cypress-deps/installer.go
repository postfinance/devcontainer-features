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
	return installer.Tools.Apt.InstallDependencies(
		"libgtk2.0-0",
		"libgtk-3-0",
		"libgbm-dev",
		"libnotify-dev",
		"libnss3",
		"libxss1",
		"libasound2",
		"libxtst6",
		"xauth",
		"xvfb",
	)
}
