package installer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/roemer/goext"
)

const (
	AMD64 = "amd64"
	ARM64 = "arm64"
)

type system struct{}

// Installs the given binary to /usr/local/bin with the given name.
func (s *system) InstallBinaryToUsrLocalBin(binaryPath string, targetBinaryName string) error {
	return s.InstallBinaryToDir(binaryPath, "/usr/local/bin", targetBinaryName)
}

// Installs the binary to the given directory with the given name.
func (s *system) InstallBinaryToDir(binaryPath string, targetDir string, targetBinaryName string) error {
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return fmt.Errorf("failed to create target directory %q: %w", targetDir, err)
	}
	return goext.CmdRunners.Console.Run("install", "-m", "0755", binaryPath, filepath.Join(targetDir, targetBinaryName))
}

// Downloads a binary from the given url and installs it to /usr/local/bin with the given name.
func (s *system) DownloadBinaryToUsrLocalBin(url string, progressName string, targetBinaryName string) error {
	return s.DownloadBinaryToDir(url, progressName, "/usr/local/bin", targetBinaryName)
}

// Downloads a binary from the given url and installs it to the given directory with the given name.
func (s *system) DownloadBinaryToDir(url string, progressName string, targetDir string, targetBinaryName string) error {
	tempDir, err := os.MkdirTemp("", "binary-download-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)
	tempPath := filepath.Join(tempDir, targetBinaryName)
	if err := Tools.Download.ToFile(url, tempPath, progressName); err != nil {
		return err
	}
	return s.InstallBinaryToDir(tempPath, targetDir, targetBinaryName)
}

func (s *system) InstallPackages(packages ...string) error {
	osInfo, err := s.GetOsInfo()
	if err != nil {
		return err
	}
	return s.InstallPackagesForOs(osInfo, packages...)
}

func (s *system) InstallPackagesForOs(osInfo *OsInfo, packages ...string) error {
	switch {
	case osInfo.IsDebian(), osInfo.IsUbuntu():
		return Tools.Apt.InstallDependencies(packages...)
	case osInfo.IsAlpine():
		return Tools.Apk.InstallDependencies(packages...)
	default:
		return fmt.Errorf("unsupported OS vendor: %s", osInfo.Vendor)
	}
}

func (s *system) InstallPackagesByOs(f func(osInfo *OsInfo) ([]string, error)) error {
	osInfo, err := s.GetOsInfo()
	if err != nil {
		return err
	}
	packages, err := f(osInfo)
	if err != nil {
		return err
	}
	if packages == nil {
		return nil
	}
	return s.InstallPackagesForOs(osInfo, packages...)
}

func (s *system) MapArchitecture(mapping map[string]string) (string, error) {
	mappedValue, ok := mapping[runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
	return mappedValue, nil
}

func (s *system) GetOsInfo() (*OsInfo, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	infoMap := map[string]string{}
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		// Remove surrounding quotes if present
		val := strings.Trim(parts[1], `"`)
		infoMap[parts[0]] = val
	}
	return &OsInfo{
		Vendor:    infoMap["ID"],
		Codename:  infoMap["VERSION_CODENAME"],
		VersionId: infoMap["VERSION_ID"],
	}, nil
}

type OsInfo struct {
	Vendor    string
	Codename  string
	VersionId string
}

func (v *OsInfo) IsDebian() bool {
	return v.Vendor == "debian"
}

func (v *OsInfo) IsUbuntu() bool {
	return v.Vendor == "ubuntu"
}

func (v *OsInfo) IsAlpine() bool {
	return v.Vendor == "alpine"
}

func (v *OsInfo) MajorVersion() int {
	var major int
	fmt.Sscanf(v.VersionId, "%d", &major)
	return major
}
