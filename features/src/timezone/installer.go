package main

import (
	"flag"
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
	fmt.Println("Configuring Timezone")

	timezone := flag.String("timezone", "Etc/UTC", "The timezone to set.")
	flag.Parse()

	fmt.Println("Setting timezone to " + *timezone)

	// Write the timezone to /etc/timezone
	if err := os.WriteFile("/etc/timezone", []byte(*timezone), os.ModePerm); err != nil {
		return err
	}

	// Symlink the zoneinfo to /etc/localtime
	newTimezone := "/usr/share/zoneinfo/" + *timezone
	target := "/etc/localtime"
	// Check that the timezone exists
	if _, err := os.Lstat(newTimezone); err != nil {
		return fmt.Errorf("timezone '%s' does not exist", newTimezone)
	}
	// Delete the original timezone symlink/file (if any)
	if _, err := os.Lstat(target); err == nil {
		if err := os.Remove(target); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}
	}
	// Create the symlink
	if err := os.Symlink(newTimezone, target); err != nil {
		return err
	}

	return nil
}
