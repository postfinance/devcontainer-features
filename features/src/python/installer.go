package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/roemer/goext"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var pythonVersionRegexp *regexp.Regexp = regexp.MustCompile(`^(?P<d1>\d+)\.(?P<d2>\d+)(?:\.(?P<d3>\d+))?(?:(?P<s4>[a-z]+)(?P<d5>\d+))?$`)

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
	downloadUrl := flag.String("downloadUrl", "", "")
	pipIndex := flag.String("pipIndex", "", "")
	pipIndexUrl := flag.String("pipIndexUrl", "", "")
	pipTrustedHost := flag.String("pipTrustedHost", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrl, "https://www.python.org/ftp/python", "python-download-url")
	installer.HandleOverride(pipIndex, "", "python-pip-index")
	installer.HandleOverride(pipIndexUrl, "", "python-pip-index-url")
	installer.HandleOverride(pipTrustedHost, "", "python-pip-trusted-host")

	// Create and process the feature
	feature := installer.NewFeature("Python", false,
		&pythonComponent{
			ComponentBase:  installer.NewComponentBase("Python", *version),
			DownloadUrl:    *downloadUrl,
			PipIndex:       *pipIndex,
			PipIndexUrl:    *pipIndexUrl,
			PipTrustedHost: *pipTrustedHost,
		},
	)
	return feature.Process()
}

//////////
// Implementation
//////////

type pythonComponent struct {
	*installer.ComponentBase
	DownloadUrl    string
	PipIndex       string
	PipIndexUrl    string
	PipTrustedHost string
}

func (c *pythonComponent) GetAllVersions() ([]*gover.Version, error) {
	allTags, err := installer.Tools.GitHub.GetTags("python", "cpython")
	if err != nil {
		return nil, err
	}
	return installer.Tools.Versioning.ParseVersionsFromList(allTags, pythonVersionRegexp, true)
}

func (c *pythonComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	name := fmt.Sprintf("Python-%s", version.Raw)
	fileName := fmt.Sprintf("%s.tgz", name)
	// The folder only contains the major, minor and sometimes the patch version
	majorMinorName := fmt.Sprintf("%d.%d", version.Major(), version.Minor())
	folderName := majorMinorName
	if version.Segments[2].IsDefined() {
		folderName += fmt.Sprintf(".%d", version.Segments[2].Number)
	}
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, folderName, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Python"); err != nil {
		return err
	}
	// Extract it
	extractDir := "./python-src"
	if err := installer.Tools.Compression.ExtractTarGz(fileName, extractDir, false); err != nil {
		return err
	}
	// Build
	if err := c.buildPython(filepath.Join(extractDir, name), version); err != nil {
		return err
	}
	// Create a symlink to the installed python version
	symLinkPath := "/usr/local/python/current"
	targetPath := fmt.Sprintf("/usr/local/python/%s", folderName)
	if err := installer.Tools.FileSystem.CreateSymLink(targetPath, symLinkPath, false); err != nil {
		return err
	}
	// Configure Pip
	pythonPath := filepath.Join(targetPath, "bin", fmt.Sprintf("%s%s", "python", majorMinorName))
	if err := c.configurePip(pythonPath); err != nil {
		return err
	}
	// Create comfortable symlink for executables
	symLinkFiles := []string{
		"python", "pip", "pydoc", "idle", "python-config",
	}
	for _, symLinkFile := range symLinkFiles {
		symTargetPath := filepath.Join(targetPath, "bin", fmt.Sprintf("%s%s", symLinkFile, majorMinorName))
		symPath := filepath.Join(targetPath, "bin", symLinkFile)
		if err := installer.Tools.FileSystem.CreateSymLink(symTargetPath, symPath, false); err != nil {
			return fmt.Errorf("failed creating symlink from '%s' to '%s': %w", symPath, symTargetPath, err)
		}
	}
	// Cleanup
	if err := os.RemoveAll(extractDir); err != nil {
		return err
	}
	if err := os.Remove(fileName); err != nil {
		return err
	}

	// Other things?
	//https://github.com/devcontainers/features/blob/main/src/python/install.sh

	return nil
}

func (c *pythonComponent) buildPython(extractDir string, version *gover.Version) error {
	// Install dependencies
	if err := c.installBuildDependencies(); err != nil {
		return err
	}
	// Create the command runner
	cmdRunner := goext.CmdRunners.Console.WithWorkingDirectory(extractDir)
	// Configure
	if err := cmdRunner.Run("./configure", fmt.Sprintf("--prefix=/usr/local/python/%s/", version.Raw)); err != nil {
		return err
	}
	// Build
	if err := cmdRunner.Run("make", "-s", fmt.Sprintf("-j%d", runtime.NumCPU())); err != nil {
		return err
	}
	// Install
	if err := cmdRunner.Run("make", "install"); err != nil {
		return err
	}
	return nil
}

func (c *pythonComponent) configurePip(pythonPath string) error {
	// Ensure pip is installed
	if err := goext.CmdRunners.Console.Run(pythonPath, goext.Cmd.SplitArgs("-m ensurepip")...); err != nil {
		return err
	}
	// Configure pip
	if c.PipIndex != "" {
		if err := goext.CmdRunners.Console.Run(pythonPath, goext.Cmd.SplitArgs("-m pip config --global set global.index", c.PipIndex)...); err != nil {
			return err
		}
	}
	if c.PipIndexUrl != "" {
		if err := goext.CmdRunners.Console.Run(pythonPath, goext.Cmd.SplitArgs("-m pip config --global set global.index-url", c.PipIndexUrl)...); err != nil {
			return err
		}
	}
	if c.PipTrustedHost != "" {
		if err := goext.CmdRunners.Console.Run(pythonPath, goext.Cmd.SplitArgs("-m pip config --global set global.trusted-host", c.PipTrustedHost)...); err != nil {
			return err
		}
	}
	return nil
}

func (c *pythonComponent) installBuildDependencies() error {
	return installer.Tools.System.InstallPackagesByOs(func(osInfo *installer.OsInfo) ([]string, error) {
		if osInfo.IsDebian() || osInfo.IsUbuntu() {
			return []string{
				"build-essential",
				"gdb",
				"lcov",
				"pkg-config",
				"libbz2-dev",
				"libffi-dev",
				"libgdbm-dev",
				"libgdbm-compat-dev",
				"liblzma-dev",
				"libncurses5-dev",
				"libreadline-dev",
				"libsqlite3-dev",
				"libssl-dev",
				"tk-dev",
				"uuid-dev",
				"zlib1g-dev",
			}, nil
		} else if osInfo.IsAlpine() {
			return []string{
				"bluez-dev",
				"bzip2-dev",
				"dpkg-dev",
				"dpkg",
				"findutils",
				"gcc",
				"gdbm-dev",
				"gnupg",
				"libc-dev",
				"libffi-dev",
				"libnsl-dev",
				"libtirpc-dev",
				"linux-headers",
				"make",
				"ncurses-dev",
				"openssl-dev",
				"pax-utils",
				"readline-dev",
				"sqlite-dev",
				"tar",
				"tcl-dev",
				"tk",
				"tk-dev",
				"util-linux-dev",
				"xz",
				"xz-dev",
				"zlib-dev",
				"zstd-dev",
			}, nil
		}
		return nil, fmt.Errorf("unsupported OS vendor: %s", osInfo.Vendor)
	})
}
