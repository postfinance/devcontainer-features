package installer

import (
	"os"

	"github.com/roemer/gotaskr/execr"
)

// Run runs an executable with the given arguments and with sudo, if needed
func Run(outputToConsole bool, executable string, arguments ...string) error {
	var args []string
	if isDevContainer() {
		// the local execution (inside dev container) needs sudo -> e.g. sudo apt-get update
		args = append([]string{executable}, arguments...)
		executable = "sudo"
	} else {
		args = arguments
	}

	if err := execr.Run(outputToConsole, executable, args...); err != nil {
		return err
	}

	return nil
}

func isDevContainer() bool {
	isDevContainer := os.Getenv("IS_DEVCONTAINER")
	if isDevContainer == "" || isDevContainer != "true" {
		return false
	}
	return true
}
