package installer

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/roemer/gotaskr/execr"
)

type fileSystem struct{}

func (f *fileSystem) EnsureLoginShellPath() error {
	restoreEnvScriptPath := "/etc/profile.d/00-restore-env.sh"

	// Remove the existing file
	err := os.Remove(restoreEnvScriptPath)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	// Get the current PATH
	originalPath := os.Getenv("PATH")

	// Get the path for a login shell
	loginPath, _, err := execr.RunGetOutput(false, "sh", "-lc", "echo $PATH")
	if err != nil {
		return err
	}

	// Replace occurrences of the login PATH with $PATH
	newPath := strings.ReplaceAll(originalPath, loginPath, "$PATH")

	// Create the content for the shell script
	scriptContent := fmt.Sprintf("export PATH=%s\n", newPath)

	// Write the content back
	err = os.WriteFile(restoreEnvScriptPath, []byte(scriptContent), 0755)
	if err != nil {
		return err
	}

	return nil
}

func (f *fileSystem) MoveFolder(src, dest string, includeFolder bool) error {
	if err := filepath.WalkDir(src, func(currPath string, d fs.DirEntry, werr error) error {
		// Handle a walk error
		if werr != nil {
			return werr
		}

		// Get the relative path
		relPath, err := filepath.Rel(filepath.Dir(src), currPath)
		if err != nil {
			return err
		}
		// Build the target path
		if !includeFolder {
			// Remove the first part of the relative path
			parts := strings.Split(relPath, string(filepath.Separator))
			relPath = strings.Join(parts[1:], string(filepath.Separator))
			if relPath == "" {
				// This is the initial folder which is now empty
				return nil
			}
		}
		targetPath := filepath.Join(dest, relPath)

		// Create the folder / move the file
		if d.IsDir() {
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed creating folder '%s': %w", targetPath, err)
			}
		} else {
			if err := os.Rename(currPath, targetPath); err != nil {
				return fmt.Errorf("failed creating file '%s': %w", targetPath, err)
			}
		}

		return nil
	}); err != nil {
		return err
	}
	return nil
}

// Moves all given folders to the destination including the folder name itself (merging if they already exist).
func (f *fileSystem) MoveFolders(src []string, dest string) error {
	for _, currSrc := range src {
		if err := f.MoveFolder(currSrc, dest, true); err != nil {
			return err
		}
	}
	return nil
}

func (f *fileSystem) MoveFile(src, dest string) error {
	dir := filepath.Dir(dest)
	os.MkdirAll(dir, os.ModePerm)
	return os.Rename(src, dest)
}

func (f *fileSystem) CreateSymLink(targetPath string, symLinkPath string, allowNotExistingTarget bool) error {
	// Check if the symlink path exists and delete it then
	if _, err := os.Lstat(symLinkPath); err == nil {
		if err := os.Remove(symLinkPath); err != nil {
			return fmt.Errorf("failed to unlink: %+v", err)
		}
	}
	// Create the symlink
	if err := os.Symlink(targetPath, symLinkPath); err != nil {
		if allowNotExistingTarget && errors.Is(err, fs.ErrNotExist) {
			// Do not give an error if the symlink points to a not existing target
			return nil
		}
		return err
	}
	return nil
}
