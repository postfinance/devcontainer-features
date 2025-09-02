package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"
	"time"

	"builder/installer"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var chromeVersionRegex *regexp.Regexp = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)\.(\d+)$`)
var firefoxVersionRegex *regexp.Regexp = regexp.MustCompile(`^(\d+)\.(\d+)(?:\.(\d+))?(?:\.(\d+))?(?:([a-z]+)(\d+))?$`)

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
	chromeVersion := flag.String("chromeVersion", "none", "")
	chromeDownloadUrl := flag.String("chromeDownloadUrl", "", "")
	chromeVersionsUrl := flag.String("chromeVersionsUrl", "", "")
	chromeTestingVersionsUrl := flag.String("chromeTestingVersionsUrl", "", "")
	useChromeForTesting := flag.Bool("useChromeForTesting", true, "")

	firefoxVersion := flag.String("firefoxVersion", "none", "")
	firefoxDownloadUrl := flag.String("firefoxDownloadUrl", "", "")
	firefoxVersionsUrl := flag.String("firefoxVersionsUrl", "", "")
	flag.Parse()

	// Load global overrides
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	// Apply override logic for URLs
	installer.HandleOverride(chromeDownloadUrl, "https://dl.google.com/linux/chrome/deb/pool/main/g/google-chrome-stable", "chrome-download-url")
	installer.HandleOverride(chromeVersionsUrl, "https://versionhistory.googleapis.com/v1/chrome/platforms/linux/channels/stable/versions", "chrome-versions-url")
	installer.HandleOverride(chromeTestingVersionsUrl, "https://googlechromelabs.github.io/chrome-for-testing/known-good-versions.json", "chrome-testing-versions-url")

	installer.HandleOverride(firefoxDownloadUrl, "https://download-installer.cdn.mozilla.net/pub/firefox/releases", "firefox-download-url")
	installer.HandleOverride(firefoxVersionsUrl, "https://product-details.mozilla.org/1.0/firefox.json", "firefox-versions-url")

	// Choose the correct chrome component
	var chromeComponentToInstall installer.IComponent
	if *useChromeForTesting {
		chromeComponentToInstall = &chromeForTestingComponent{
			ComponentBase:         installer.NewComponentBase("Chrome Test", *chromeVersion),
			TestingVersionsUrl:    *chromeTestingVersionsUrl,
			ChromeDownloadBaseUrl: *chromeDownloadUrl,
		}
	} else {
		chromeComponentToInstall = &chromeComponent{
			ComponentBase:         installer.NewComponentBase("Chrome", *chromeVersion),
			VersionsUrl:           *chromeVersionsUrl,
			ChromeDownloadBaseUrl: *chromeDownloadUrl,
		}
	}

	// Create and process the feature
	feature := installer.NewFeature("Browsers", false,
		chromeComponentToInstall,
		&firefoxComponent{
			ComponentBase:          installer.NewComponentBase("Firefox", *firefoxVersion),
			VersionsUrl:            *firefoxVersionsUrl,
			FirefoxDownloadBaseUrl: *firefoxDownloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type chromeComponent struct {
	*installer.ComponentBase
	VersionsUrl           string
	ChromeDownloadBaseUrl string
}

func (c *chromeComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(c.VersionsUrl)
	if err != nil {
		return nil, err
	}
	var jsonData struct {
		Versions []struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"versions"`
		NextPageToken string `json:"nextPageToken"`
	}
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	for _, entry := range jsonData.Versions {
		version := gover.MustParseVersionFromRegex(entry.Version, chromeVersionRegex)
		allVersions = append(allVersions, version)
	}
	return allVersions, nil
}

func (c *chromeComponent) InstallVersion(version *gover.Version) error {
	// Download the file
	fileName := "google-chrome-stable_amd64.deb"
	downloadUrl := fmt.Sprintf("%s/google-chrome-stable_%s-1_amd64.deb", c.ChromeDownloadBaseUrl, version.Raw)
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Chrome"); err != nil {
		return err
	}
	// Install it

	if err := installer.Tools.Apt.InstallLocalPackage(fileName); err != nil {
		return err
	}
	// Cleanup
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}
	return nil
}

type chromeForTestingComponent struct {
	*installer.ComponentBase
	TestingVersionsUrl    string
	ChromeDownloadBaseUrl string
}

func (c *chromeForTestingComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(c.TestingVersionsUrl)
	if err != nil {
		return nil, err
	}
	var jsonData struct {
		Timestamp time.Time `json:"timestamp"`
		Versions  []struct {
			Version  string `json:"version"`
			Revision string `json:"revision"`
		} `json:"versions"`
	}
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	for _, entry := range jsonData.Versions {
		version := gover.MustParseVersionFromRegex(entry.Version, chromeVersionRegex)
		allVersions = append(allVersions, version)
	}
	return allVersions, nil
}

func (c *chromeForTestingComponent) InstallVersion(version *gover.Version) error {
	// Remove known-good-versions.json from end if present using split
	downloadBaseUrl := c.TestingVersionsUrl
	if len(downloadBaseUrl) > 0 {
		parts := regexp.MustCompile(`known-good-versions\.json$`).Split(downloadBaseUrl, -1)
		if len(parts) > 0 {
			downloadBaseUrl = parts[0]
		}
	}
	metaFileContent, err := installer.Tools.Download.AsBytes(fmt.Sprintf("%s/%s.json", downloadBaseUrl, version.Raw))
	if err != nil {
		return err
	}
	var jsonMetaData struct {
		Version   string `json:"version"`
		Revision  string `json:"revision"`
		Downloads struct {
			Chrome []struct {
				Platform string `json:"platform"`
				URL      string `json:"url"`
			} `json:"chrome"`
		} `json:"downloads"`
	}
	if err := json.Unmarshal(metaFileContent, &jsonMetaData); err != nil {
		return err
	}
	downloadUrl := ""
	for _, dl := range jsonMetaData.Downloads.Chrome {
		if dl.Platform == "linux64" {
			downloadUrl = dl.URL
			break
		}
	}
	if downloadUrl == "" {
		return fmt.Errorf("no download url for linux found")
	}

	// Download the file
	fileName := "chrome.zip"
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Chrome Test"); err != nil {
		return err
	}
	// Extract it
	if err := installer.Tools.Compression.ExtractZip(fileName, "/usr/local/chrome", true); err != nil {
		return err
	}
	// Create symlink to binary
	if err := installer.Tools.FileSystem.CreateSymLink("/usr/local/chrome/chrome", "/usr/local/bin/chrome", false); err != nil {
		return err
	}

	// Install dependencies
	// Taken from https://source.chromium.org/chromium/chromium/src/+/main:chrome/installer/linux/debian/dist_package_versions.json
	if err := installer.Apt.InstallDependencies(
		"libasound2",
		"libatk-bridge2.0-0",
		"libatk1.0-0",
		"libatspi2.0-0",
		"libc6",
		"libcairo2",
		"libcups2",
		"libdbus-1-3",
		"libdrm2",
		"libexpat1",
		"libgbm1",
		"libgcc-s1",
		"libglib2.0-0",
		"libnspr4",
		"libnss3",
		"libpango-1.0-0",
		"libpangocairo-1.0-0",
		"libstdc++6",
		"libuuid1",
		"libx11-6",
		"libx11-xcb1",
		"libxcb-dri3-0",
		"libxcb1",
		"libxcomposite1",
		"libxcursor1",
		"libxdamage1",
		"libxext6",
		"libxfixes3",
		"libxi6",
		"libxkbcommon0",
		"libxrandr2",
		"libxrender1",
		"libxshmfence1",
		"libxss1",
		"libxtst6",
	); err != nil {
		return err
	}

	// Cleanup
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}

	return nil
}

type firefoxComponent struct {
	*installer.ComponentBase
	VersionsUrl            string
	FirefoxDownloadBaseUrl string
}

func (c *firefoxComponent) GetAllVersions() ([]*gover.Version, error) {
	versionFileContent, err := installer.Tools.Download.AsBytes(c.VersionsUrl)
	if err != nil {
		return nil, err
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(versionFileContent, &jsonData); err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	for _, entry := range jsonData["releases"].(map[string]interface{}) {
		versionString := entry.(map[string]interface{})["version"].(string)
		version := gover.MustParseVersionFromRegex(versionString, firefoxVersionRegex)
		allVersions = append(allVersions, version)
	}
	return allVersions, nil
}

func (c *firefoxComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x86_64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}
	fileName := "firefox.tar.xz"
	downloadUrl := fmt.Sprintf("%s/%s/linux-%s/en-US/firefox-%s.tar.xz", c.FirefoxDownloadBaseUrl, version.Raw, archPart, version.Raw)
	if version.LessThan(gover.ParseSimple("135")) {
		fileName = "firefox.tar.bz2"
		downloadUrl = fmt.Sprintf("%s/%s/linux-%s/en-US/firefox-%s.tar.bz2", c.FirefoxDownloadBaseUrl, version.Raw, archPart, version.Raw)
	}
	// Download the file
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Firefox"); err != nil {
		return err
	}

	// Extract it
	if err := installer.Tools.Compression.Extract(fileName, "/usr/local/firefox", true); err != nil {
		return err
	}

	// Create symlink to binary
	if err := installer.Tools.FileSystem.CreateSymLink("/usr/local/firefox/firefox", "/usr/local/bin/firefox", false); err != nil {
		return err
	}

	// Install dependencies
	if err := installer.Apt.InstallDependencies(
		"libasound2",
		"libgtk-3-0",
		"libx11-xcb1",
		"libxml2-utils",
		"libxt6",
		"libxtst6",
	); err != nil {
		return err
	}

	// Cleanup
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}

	return nil
}
