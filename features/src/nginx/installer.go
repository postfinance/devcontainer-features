package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var versionRegexp *regexp.Regexp = regexp.MustCompile(`^(\d+).(\d+).(\d+)-(\d+)$`)
var indexLineRegexp *regexp.Regexp = regexp.MustCompile(`^.*<a href="(nginx_([0-9\.\-]+)~([a-z]+)_amd64.deb)">.*$`)

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
	stableOnly := flag.Bool("stableOnly", false, "")
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file (global/per-feature overrides)
	if err := installer.LoadOverrides(); err != nil {
		return err
	}
	// Handle overrides and defaults
	installer.HandleOverride(downloadUrl, "https://nginx.org", "nginx-download-url")

	// Fetch the os info
	osInfo, err := installer.Tools.System.GetOsInfo()
	if err != nil {
		return fmt.Errorf("failed getting OS info: %w", err)
	}

	// Create and process the feature
	feature := installer.NewFeature("Nginx", false,
		&nginxComponent{
			ComponentBase: installer.NewComponentBase("Nginx", *version),
			stableOnly:    *stableOnly,
			osInfo:        osInfo,
			downloadUrl:   *downloadUrl,
		},
	)
	return feature.Process()
}

//////////
// Implementation
//////////

type nginxComponent struct {
	*installer.ComponentBase
	stableOnly  bool
	osInfo      *installer.OsInfo
	downloadUrl string
}

func (c *nginxComponent) IsFullVersion(referenceVersion *gover.Version) bool {
	return versionRegexp.Match([]byte(referenceVersion.Raw))
}

func (c *nginxComponent) GetAllVersions() ([]*gover.Version, error) {
	versions := []*gover.Version{}
	// Stable
	stableVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndexFunc(c.getStableUrl(), c.lineExtractFunc)
	if err != nil {
		return nil, err
	}
	versions = append(versions, stableVersions...)
	// Mainline
	if !c.stableOnly {
		mainlineVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndexFunc(c.getMainlineUrl(), c.lineExtractFunc)
		if err != nil {
			return nil, err
		}
		versions = append(versions, mainlineVersions...)
	}
	return versions, nil
}

func (c *nginxComponent) InstallVersion(version *gover.Version) error {
	fileName := "nginx.deb"
	urls := []string{}
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "amd64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return err
	}
	debName := fmt.Sprintf("nginx_%s~%s_%s.deb", version.Raw, c.osInfo.Codename, archPart)
	urls = append(urls, c.getStableUrl()+debName)
	if !c.stableOnly {
		urls = append(urls, c.getMainlineUrl()+debName)
	}
	// Download the file with fallback urls
	var lastErr error
	for _, downloadUrl := range urls {
		if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Nginx"); err == nil {
			// Download succeeded, proceed with install
			lastErr = nil
			break
		} else {
			lastErr = err
		}
	}
	if lastErr != nil {
		return lastErr
	}
	if err := installer.Tools.Apt.InstallLocalPackage(fileName); err != nil {
		return err
	}
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}
	return nil

}

func (c *nginxComponent) getStableUrl() string {
	return fmt.Sprintf("%s/packages/%s/pool/nginx/n/nginx/", c.downloadUrl, c.osInfo.Vendor)
}

func (c *nginxComponent) getMainlineUrl() string {
	return fmt.Sprintf("%s/packages/mainline/%s/pool/nginx/n/nginx/", c.downloadUrl, c.osInfo.Vendor)
}

func (c *nginxComponent) lineExtractFunc(url, line string) (*gover.Version, error) {
	if match := indexLineRegexp.FindStringSubmatch(line); match != nil {
		versionString := match[2]
		codename := match[3]
		if codename == c.osInfo.Codename {
			version := gover.MustParseVersionFromRegex(versionString, versionRegexp)
			return version, nil
		}
	}
	return nil, nil
}
