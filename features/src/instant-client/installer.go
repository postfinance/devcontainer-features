package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/roemer/gotaskr/execr"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

var versionRegexp *regexp.Regexp = regexp.MustCompile(`^(\d+).(\d+).(\d+).(\d+).(\d+)$`)
var indexLineRegexp *regexp.Regexp = regexp.MustCompile(`^.*<a.*href='.*download\.oracle\.com(.*instantclient-basic-.*-(\d+(?:\.\d+){4})\w*\.zip)'.*$`)

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
	version := flag.String("version", "latest", "The version of Instant Client to install.")
	versionsUrl := flag.String("versionsUrl", "", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	flag.Parse()

	// Load settings from an external file (global/per-feature overrides)
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(versionsUrl, "https://www.oracle.com/database/technologies/instant-client/linux-x86-64-downloads.html", "instant-client-versions-url")
	installer.HandleOverride(downloadUrl, "https://download.oracle.com", "instant-client-download-url")

	// Create and process the feature
	feature := installer.NewFeature("Oracle Instant Client", false,
		&instantClientComponent{
			ComponentBase: installer.NewComponentBase("Basic Package", *version),
			versionsUrl:   *versionsUrl,
			downloadUrl:   *downloadUrl,
		})
	return feature.Process()
}

//////////
// Implementation
//////////

type instantClientComponent struct {
	*installer.ComponentBase
	versionsUrl string
	downloadUrl string
}

func (c *instantClientComponent) IsFullVersion(referenceVersion *gover.Version) bool {
	return len(referenceVersion.Segments) == 5 && referenceVersion.DefinedSegmentCount() == 5
}

func (c *instantClientComponent) createDownloadURLForVersion(version *gover.Version) (string, error) {
	zipVersion := version.Raw
	// Versions below 23 have a dbru suffix
	if version.Major() < 23 {
		zipVersion = fmt.Sprintf("%sdbru", version.Raw)
	}
	majorMinor := fmt.Sprintf("%d%d", version.Major(), version.Minor())
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s/otn_software/linux/instantclient/%s/instantclient-basic-linux.%s-%s.zip",
		c.downloadUrl,
		fmt.Sprintf("%s%s", majorMinor, strings.Repeat("0", 7-len(majorMinor))),
		archPart,
		zipVersion,
	), nil
}

func (c *instantClientComponent) GetAllVersions() ([]*gover.Version, error) {
	// Parse the latest versions from download page
	versions := []*gover.Version{}
	stableVersions, err := installer.Tools.Http.GetVersionsFromHtmlIndexFunc(c.versionsUrl, c.lineExtractFunc)
	if err != nil {
		return nil, err
	}
	versions = append(versions, stableVersions...)
	return versions, nil
}

func (c *instantClientComponent) InstallVersion(version *gover.Version) error {
	fileName := "instant-client.zip"
	downloadUrl, err := c.createDownloadURLForVersion(version)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Instant Client"); err != nil {
		return err
	}
	rootFolder, err := installer.Tools.Compression.GetRootFolderFromZip(fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Compression.ExtractZip(fileName, "/opt/oracle", false); err != nil {
		return err
	}
	if err := installer.Tools.Apt.InstallDependencies("libaio1"); err != nil {
		return err
	}
	if err := os.WriteFile("/etc/ld.so.conf.d/oracle-instantclient.conf", []byte(path.Join("/opt/oracle", rootFolder)), 0644); err != nil {
		return err
	}
	if err := execr.Run(true, "ldconfig"); err != nil {
		return err
	}
	// Cleanup
	if err := os.RemoveAll(fileName); err != nil {
		return err
	}
	return nil
}

func (c *instantClientComponent) lineExtractFunc(url, line string) (*gover.Version, error) {
	if match := indexLineRegexp.FindStringSubmatch(line); match != nil {
		fullName := match[2]

		version := gover.MustParseVersionFromRegex(fullName, versionRegexp)
		return version, nil
	}
	return nil, nil
}
