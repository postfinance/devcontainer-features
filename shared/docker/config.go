package docker

import (
	"builder/installer"
	"os"
	"os/user"
	"path/filepath"
	"strconv"
)

func InstallClientConfig(configPath string) error {
	if configPath != "" {
		fileContent, err := installer.ReadFileFromUrlOrLocal(configPath)
		if err != nil {
			return err
		}
		userName := os.Getenv("_REMOTE_USER")
		homeDir := os.Getenv("_REMOTE_USER_HOME")
		if homeDir == "" {
			homeDir = filepath.Join("/home", userName)
		}
		dockerDir := filepath.Join(homeDir, ".docker")
		configDest := filepath.Join(dockerDir, "config.json")
		// Ensure directory exists
		if err := os.MkdirAll(dockerDir, 0700); err != nil {
			return err
		}
		// Write config file
		if err := os.WriteFile(configDest, fileContent, 0600); err != nil {
			return err
		}
		// Set ownership
		usr, err := user.Lookup(userName)
		if err != nil {
			return err
		}
		uid, _ := strconv.Atoi(usr.Uid)
		gid, _ := strconv.Atoi(usr.Gid)
		if err := os.Chown(dockerDir, uid, gid); err != nil {
			return err
		}
		if err := os.Chown(configDest, uid, gid); err != nil {
			return err
		}
	}
	return nil
}
