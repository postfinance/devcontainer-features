package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

type Product int

const (
	sdk Product = iota
	runtime
	aspNetRuntime
)

func (c Product) String() string {
	switch c {
	case sdk:
		return "Sdk"
	case runtime:
		return "Runtime"
	case aspNetRuntime:
		return "aspnetcore/Runtime"
	default:
		return ""
	}
}

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
	additionalVersions := flag.String("additionalVersions", "", "")
	dotnetRuntimeVersions := flag.String("dotnetRuntimeVersions", "", "")
	aspNetCoreRuntimeVersions := flag.String("aspNetCoreRuntimeVersions", "", "")
	workloads := flag.String("workloads", "", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	versionsUrl := flag.String("versionsUrl", "", "")
	nugetConfigPath := flag.String("nugetConfigPath", "", "")

	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrl, "https://builds.dotnet.microsoft.com/dotnet", "dotnet-download-url")
	installer.HandleOverride(versionsUrl, "https://builds.dotnet.microsoft.com/dotnet", "dotnet-versions-url")
	installer.HandleOverride(nugetConfigPath, "", "dotnet-nuget-config-path")

	// Handle multi value fields
	var allSdks = []string{*version}
	if len(*additionalVersions) > 0 {
		allSdks = append(allSdks, strings.Split(*additionalVersions, ",")...)
	}
	var additionalRuntimes = []string{}
	if len(*dotnetRuntimeVersions) > 0 {
		additionalRuntimes = strings.Split(*dotnetRuntimeVersions, ",")
	}
	var additionalaspNetCoreRuntimes = []string{}
	if len(*aspNetCoreRuntimeVersions) > 0 {
		additionalaspNetCoreRuntimes = strings.Split(*aspNetCoreRuntimeVersions, ",")
	}
	// Create the feature
	feature := installer.NewFeature(".NET", true)
	if *nugetConfigPath != "" {
		feature.AddComponents(&nugetConfigComponent{
			ComponentBase:   installer.NewComponentBase("Nuget Config", installer.VERSION_IRRELEVANT),
			NugetConfigPath: *nugetConfigPath,
		})
	}
	// add sdks
	for _, sdkVersion := range allSdks {
		component := &sdkComponent{
			ComponentBase: installer.NewComponentBase(fmt.Sprintf("SDK [%s]", sdkVersion), strings.TrimSpace(sdkVersion)),
			DownloadUrl:   *downloadUrl,
			VersionsUrl:   *versionsUrl,
		}
		feature.AddComponents(component)
	}
	// add runtimes
	for _, runtimeVersion := range additionalRuntimes {
		component := &runtimeComponent{
			ComponentBase: installer.NewComponentBase(fmt.Sprintf("Runtime [%s]", runtimeVersion), strings.TrimSpace(runtimeVersion)),
			DownloadUrl:   *downloadUrl,
			VersionsUrl:   *versionsUrl,
		}
		feature.AddComponents(component)
	}
	// add runtimes
	for _, aspCoreVersion := range additionalaspNetCoreRuntimes {
		component := &aspNetRuntimeComponent{
			ComponentBase: installer.NewComponentBase(fmt.Sprintf("ASP.NET Core runtime [%s]", aspCoreVersion), strings.TrimSpace(aspCoreVersion)),
			DownloadUrl:   *downloadUrl,
			VersionsUrl:   *versionsUrl,
		}
		feature.AddComponents(component)
	}
	// workloads
	if len(*workloads) > 0 {
		feature.AddComponents(&workloadComponent{
			ComponentBase: installer.NewComponentBase("Workloads", installer.VERSION_IRRELEVANT),
			workloads:     strings.Split(strings.ReplaceAll(*workloads, " ", ""), ","),
		})
	}
	// Component to create a symlink, but only if an sdk was installed
	if len(allSdks) > 0 {
		feature.AddComponents(&symlinkComponent{
			ComponentBase: installer.NewComponentBase("SymLink", installer.VERSION_IRRELEVANT),
		})
	}
	// Process the feature
	return feature.Process()
}

//////////
// Implementation
//////////

type sdkComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
}

func (c *sdkComponent) GetAllVersions() ([]*gover.Version, error) {
	latestVersion, err := resolveSdkVersion(c.VersionsUrl, c.GetRequestedVersion())
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, gover.RegexpSemver)
	if err != nil {
		return nil, err
	}
	return []*gover.Version{version}, err
}

func (c *sdkComponent) GetLatestVersion() (*gover.Version, error) {
	latestVersion, err := resolveSdkVersion(c.VersionsUrl, installer.VERSION_LTS)
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, gover.RegexpSemver)
	if err != nil {
		return nil, err
	}
	return version, err
}

func (c *sdkComponent) InstallVersion(version *gover.Version) error {
	return installSdk(c.DownloadUrl, version)
}

type runtimeComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
}

func (c *runtimeComponent) GetAllVersions() ([]*gover.Version, error) {
	latestVersion, err := resolveRuntimeVersion(c.VersionsUrl, c.GetRequestedVersion())
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, gover.RegexpSemver)
	if err != nil {
		return nil, err
	}
	return []*gover.Version{version}, err
}

func (c *runtimeComponent) GetLatestVersion() (*gover.Version, error) {
	latestVersion, err := resolveRuntimeVersion(c.VersionsUrl, installer.VERSION_LTS)
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, gover.RegexpSemver)
	if err != nil {
		return nil, err
	}
	return version, err
}

func (c *runtimeComponent) InstallVersion(version *gover.Version) error {
	return installRuntime(c.DownloadUrl, version)
}

type aspNetRuntimeComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
}

func (c *aspNetRuntimeComponent) GetAllVersions() ([]*gover.Version, error) {
	latestVersion, err := resolveAspNetRuntimeVersion(c.VersionsUrl, c.GetRequestedVersion())
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, gover.RegexpSemver)
	if err != nil {
		return nil, err
	}
	return []*gover.Version{version}, err
}

func (c *aspNetRuntimeComponent) GetLatestVersion() (*gover.Version, error) {
	latestVersion, err := resolveAspNetRuntimeVersion(c.VersionsUrl, installer.VERSION_LTS)
	if err != nil {
		return nil, err
	}
	version, err := gover.ParseVersionFromRegex(latestVersion, gover.RegexpSemver)
	if err != nil {
		return nil, err
	}
	return version, err
}

func (c *aspNetRuntimeComponent) InstallVersion(version *gover.Version) error {
	return installAspNetRuntime(c.DownloadUrl, version)
}

func installSdk(downloadUrl string, version *gover.Version) error {
	return installDotnetBinary(downloadUrl, sdk, "dotnet-sdk", "downloading sdk", version)
}

func installRuntime(downloadUrl string, version *gover.Version) error {
	return installDotnetBinary(downloadUrl, runtime, "dotnet-runtime", "downloading runtime", version)
}

func installAspNetRuntime(downloadUrl string, version *gover.Version) error {
	return installDotnetBinary(downloadUrl, aspNetRuntime, "aspnetcore-runtime", "downloading ASP.NET runtime", version)
}

func installDotnetBinary(downloadUrl string, product Product, fileName string, progressName string, version *gover.Version) error {
	osInfo, err := installer.Tools.System.GetOsInfo()
	if err != nil {
		return err
	}

	// Determine the architecture part of the url
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	arch := ""
	if osInfo.IsAlpine() {
		arch = fmt.Sprintf("linux-musl-%s", archPart)
		installer.Tools.System.InstallPackages("ca-certificates", "libgcc", "libssl3", "libstdc++", "zlib", "icu-libs", "icu-data-full", "tzdata", "krb5")
	} else if osInfo.IsDebian() || osInfo.IsUbuntu() {
		arch = fmt.Sprintf("linux-%s", archPart)
	} else {
		return fmt.Errorf("unsupported OS for .NET install")
	}

	// Download file
	downloadedFileName := fmt.Sprintf("%s-%s-%s.tar.gz", fileName, version.Raw, arch)
	fullUrl := fmt.Sprintf("%s/%s/%s/%s", downloadUrl, product, version.Raw, downloadedFileName)
	if err := installer.Tools.Download.ToFile(fullUrl, downloadedFileName, progressName); err != nil {
		return err
	}

	// Extract it
	if err := installer.Tools.Compression.ExtractTarGz(downloadedFileName, os.Getenv("DOTNET_ROOT"), false); err != nil {
		return err
	}
	// Cleanup
	if err := os.Remove(downloadedFileName); err != nil {
		return err
	}
	return nil
}

func resolveSdkVersion(versionsUrl string, requestedVersion string) (string, error) {
	return resolveVersion(versionsUrl, requestedVersion, sdk)
}

func resolveRuntimeVersion(versionsUrl string, requestedVersion string) (string, error) {
	return resolveVersion(versionsUrl, requestedVersion, runtime)
}

func resolveAspNetRuntimeVersion(versionsUrl string, requestedVersion string) (string, error) {
	return resolveVersion(versionsUrl, requestedVersion, aspNetRuntime)
}

func resolveVersion(versionsUrl string, requestedVersion string, product Product) (string, error) {
	var latestVersion string
	regexChannel := regexp.MustCompile(`^(sts|lts|\d+\.\d+)$`)
	if regexChannel.MatchString(strings.ToLower(requestedVersion)) {
		// 4.0 works as channel
		// 8.0.1xx feature band should work according to docs but does somehow work only for some versions, e.g. 8.0.2xx does not work...
		latestVersionUrl := fmt.Sprintf("%s/%s/%s/latest.version", versionsUrl, product, strings.ToUpper(requestedVersion))

		version, err := installer.Tools.Download.AsString(latestVersionUrl)
		if err != nil {
			return "", err
		}
		latestVersion = version
	} else {
		latestVersion = requestedVersion
	}
	return latestVersion, nil
}

type workloadComponent struct {
	*installer.ComponentBase
	workloads []string
}

func (c *workloadComponent) InstallVersion(version *gover.Version) error {
	arguments := append([]string{"workload", "install", "--temp-dir", "/tmp/dotnet-workload-temp-dir"}, c.workloads...)
	if err := execr.Run(true, "dotnet", arguments...); err != nil {
		return err
	}
	// # Clean up
	return os.RemoveAll("/tmp/dotnet-workload-temp-dir")
}

type symlinkComponent struct {
	*installer.ComponentBase
}

func (c *symlinkComponent) InstallVersion(version *gover.Version) error {
	return installer.Tools.FileSystem.CreateSymLink(fmt.Sprintf("%s/dotnet", os.Getenv("DOTNET_ROOT")), "/usr/bin/dotnet", false)
}

type nugetConfigComponent struct {
	*installer.ComponentBase
	NugetConfigPath string
}

func (c *nugetConfigComponent) InstallVersion(version *gover.Version) error {
	fileContent, err := installer.ReadFileFromUrlOrLocal(c.NugetConfigPath)
	if err != nil {
		return err
	}
	// ensure nuget.org source is disabled and we only access sources defined in the provided config file
	execr.Run(true, "dotnet", "nuget", "disable", "source", "nuget.org")
	if err := os.MkdirAll("/etc/opt/NuGet", 0755); err != nil {
		return err
	}
	return os.WriteFile("/etc/opt/NuGet/NuGetDefaults.config", fileContent, 0644)
}
