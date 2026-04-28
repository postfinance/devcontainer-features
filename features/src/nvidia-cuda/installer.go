package main

import (
	"builder/installer"
	"sync"

	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/roemer/goext"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

const nvidiaDownloadBaseURL = "https://developer.download.nvidia.com"

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
	version := flag.String("version", "latest", "")
	keyringVersion := flag.String("keyringVersion", "", "")
	installLibraries := flag.Bool("installLibraries", true, "")
	installDevLibraries := flag.Bool("installDevLibraries", true, "")
	installCompiler := flag.Bool("installCompiler", true, "")
	installTools := flag.Bool("installTools", true, "")
	additionalCudaPackages := flag.String("additionalCudaPackages", "", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrl, nvidiaDownloadBaseURL, "nvidia-cuda-download-url")

	// Create and process the feature
	feature := installer.NewFeature("NVIDIA CUDA", false,
		&cudaKeyringComponent{
			ComponentBase: installer.NewComponentBase("Keyring", *keyringVersion),
			DownloadUrl:   *downloadUrl,
		},
	)
	if *installLibraries {
		feature.AddComponents(&cudaPackageComponent{
			ComponentBase: installer.NewComponentBase("Libraries", *version),
			PackageName:   "cuda-libraries",
		})
	}
	if *installDevLibraries {
		feature.AddComponents(&cudaPackageComponent{
			ComponentBase: installer.NewComponentBase("Dev Libraries", *version),
			PackageName:   "cuda-libraries-dev",
		})
	}
	if *installCompiler {
		feature.AddComponents(&cudaPackageComponent{
			ComponentBase: installer.NewComponentBase("Compiler", *version),
			PackageName:   "cuda-compiler",
		})
	}
	if *installTools {
		feature.AddComponents(&cudaPackageComponent{
			ComponentBase: installer.NewComponentBase("Tools", *version),
			PackageName:   "cuda-tools",
		})
	}
	for packageName := range strings.SplitSeq(*additionalCudaPackages, ",") {
		packageName = strings.TrimSpace(packageName)
		if len(packageName) > 0 {
			feature.AddComponents(&cudaPackageComponent{
				ComponentBase: installer.NewComponentBase(packageName, *version),
				PackageName:   packageName,
			})
		}
	}
	return feature.Process()
}

//////////
// Implementation
//////////

type cudaKeyringComponent struct {
	*installer.ComponentBase
	DownloadUrl string
}

func (c *cudaKeyringComponent) getCudaRepo(baseUrl string) (string, error) {
	osInfo, err := installer.Tools.System.GetOsInfo()
	if err != nil {
		return "", err
	}
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return "", err
	}

	// Note: For some reason, the website started to show errors when accessing the index without a trailing slash
	if osInfo.IsDebian() {
		if archPart == "arm64" {
			return "", fmt.Errorf("No CUDA binaries are available for ARM64")
		}
		return fmt.Sprintf("%s/compute/cuda/repos/%s%d/%s/", baseUrl, osInfo.Vendor, osInfo.MajorVersion(), archPart), nil
	}
	if osInfo.IsUbuntu() {
		return fmt.Sprintf("%s/compute/cuda/repos/%s%s/%s/", baseUrl, osInfo.Vendor, strings.ReplaceAll(osInfo.VersionId, ".", ""), archPart), nil
	}
	return "", fmt.Errorf("unsupported OS: %s", osInfo.Vendor)
}

func (c *cudaKeyringComponent) GetAllVersions() ([]*gover.Version, error) {
	cudaRepo, err := c.getCudaRepo(nvidiaDownloadBaseURL)
	if err != nil {
		return nil, err
	}
	indexUrl := fmt.Sprintf("%s%s", cudaRepo, "index.html")

	allVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndex(
		indexUrl,
		regexp.MustCompile(`^.*<a href=["']cuda-keyring_([0-9\.-]+)_.*\.deb["']>.*$`),
		regexp.MustCompile(`^(\d+)\.(\d+)-(\d+)$`))
	if err != nil {
		return nil, err
	}
	return allVersions, nil
}

func (c *cudaKeyringComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	cudaRepo, err := c.getCudaRepo(c.DownloadUrl)
	if err != nil {
		return err
	}
	downloadUrl := fmt.Sprintf("%s/cuda-keyring_%s_all.deb", cudaRepo, version.Raw)
	fileName := "cuda-keyring.deb"
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "keyring"); err != nil {
		return err
	}
	// Install it
	if err := installer.Tools.Apt.InstallLocalPackage(fileName); err != nil {
		return fmt.Errorf("failed to install local package %s: %w", fileName, err)
	}
	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}

type cudaPackageComponent struct {
	*installer.ComponentBase
	PackageName string
}

func (c *cudaPackageComponent) GetAllVersions() ([]*gover.Version, error) {
	return getAllVersions(c.PackageName)
}

func (c *cudaPackageComponent) InstallVersion(version *gover.Version) error {
	packageName := fmt.Sprintf("%s-%d-%d", c.PackageName, version.Major(), version.Minor())
	return installer.Tools.System.InstallPackages(fmt.Sprintf("%s=%s", packageName, version.Raw))
}

//////////
// Internal
//////////

func getAllVersions(libraryName string) ([]*gover.Version, error) {
	// Prepare
	nameRegex := regexp.MustCompile(regexp.QuoteMeta(libraryName) + `-[0-9\-]+/`)
	versionRegex := regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)-(\d+)$`)

	// Update apt packages only once to improve performance, as multiple components will call this function
	var once sync.Once
	var updateErr error
	once.Do(func() {
		updateErr = installer.Tools.Apt.Update()
	})
	if updateErr != nil {
		return nil, updateErr
	}

	// Get the package versions
	stdout, _, err := goext.CmdRunners.Default.RunGetOutput("apt", "list", "-a", libraryName+"-*")
	if err != nil {
		return nil, err
	}

	versions := []*gover.Version{}
	for i, line := range strings.Split(stdout, "\n") {
		if i == 0 || len(line) == 0 {
			continue
		}
		if !nameRegex.MatchString(line) {
			continue
		}
		lineParts := strings.Fields(line)
		if len(lineParts) < 2 {
			// Unexpected format, skip this line to avoid panicking on out-of-range access.
			continue
		}
		version, err := gover.ParseVersionFromRegex(lineParts[1], versionRegex)
		if err != nil {
			// Version does not match the expected pattern; skip this entry safely.
			continue
		}
		versions = goext.SliceAppendIfMissing(versions, version)
	}

	return versions, nil
}
