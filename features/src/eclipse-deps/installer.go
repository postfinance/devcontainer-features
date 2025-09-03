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
	return installer.Tools.Apt.InstallDependencies(
		"libswt-gtk-4-jni",
		"libwebkit2gtk-4.0-37",
	)
}
