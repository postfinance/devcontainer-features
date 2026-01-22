package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
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
	// Handle the flags
	packages := flag.String("packages", "", "")
	flag.Parse()

	if _, err := exec.LookPath("npm"); err != nil {
		return fmt.Errorf("could not find npm, make sure to install it first.")
	}

	packagesList := strings.Split(*packages, ",")
	var pkgs []string
	for _, p := range packagesList {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			pkgs = append(pkgs, trimmed)
		}
	}
	if len(pkgs) == 0 {
		fmt.Println("No packages specified for installation.")
		return nil
	}
	for _, pkg := range pkgs {
		fmt.Printf("Installing '%s'\n", pkg)
		if err := installer.Tools.Npm.InstallGlobalPackage(pkg); err != nil {
			return err
		}
	}
	return nil
}
