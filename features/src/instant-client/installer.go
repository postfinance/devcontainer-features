package main

import (
	"builder/installer"
	"flag"
	"fmt"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/roemer/goext"
	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

const (
	defaultDownloadUrl      = "https://download.oracle.com"
	defaultVersionsUrlAMD64 = "https://www.oracle.com/database/technologies/instant-client/linux-x86-64-downloads.html"
	defaultVersionsUrlArm64 = "https://www.oracle.com/database/technologies/instant-client/linux-arm-aarch64-downloads.html"
)

var versionRegexp *regexp.Regexp = regexp.MustCompile(`^(\d+).(\d+).(\d+).(\d+).(\d+)$`)
var indexLineRegexp *regexp.Regexp = regexp.MustCompile(`<a[^>]*href=['"]([^'"]*download\.oracle\.com[^'"]*instantclient-basic-[^'"]*-(\d+(?:\.\d+){4})\w*\.zip)['"]`)

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

	defaultVersionsUrl, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: defaultVersionsUrlAMD64,
		installer.ARM64: defaultVersionsUrlArm64,
	})
	if err != nil {
		return err
	}

	installer.HandleOverride(versionsUrl, defaultVersionsUrl, "instant-client-versions-url")
	installer.HandleOverride(downloadUrl, defaultDownloadUrl, "instant-client-download-url")

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
	// Download URLs for Instant Client versions are structured as follows:
	// <baseurl>/otn_software/linux/instantclient/<subfolder>/instantclient-basic-linux.<arch>-<version><suffix>.zip
	// Examples:
	// * https://download.oracle.com/otn_software/linux/instantclient/2326300/instantclient-basic-linux.x64-23.26.3.0.0.zip
	// * https://download.oracle.com/otn_software/linux/instantclient/2380000/instantclient-basic-linux.x64-23.8.0.25.04.zip
	// * https://download.oracle.com/otn_software/linux/instantclient/199000/instantclient-basic-linux.x64-19.9.0.0.0dbru.zip
	urlPattern := "%s/otn_software/linux/instantclient/%s/instantclient-basic-linux.%s-%s.zip"

	subFolder := ""
	if version.Major() == 23 && version.Minor() < 26 {
		// For versions 23.x.x.x.x before 23.26.x.x.x the subfolder is exactly 7 digits long, with
		// the last digits being 0s and skipping all numbers after the minor version.
		const subFolderLength = 7
		majorMinor := fmt.Sprintf("%d%d", version.Major(), version.Minor())
		subFolder = fmt.Sprintf("%s%s", majorMinor, strings.Repeat("0", subFolderLength-len(majorMinor)))
	} else {
		// For all other versions, the subfolder is the version number without dots.
		subFolder = strings.ReplaceAll(version.Raw, ".", "")
	}

	zipVersion := version.Raw
	// Versions below 23 have a dbru suffix
	if version.Major() < 23 {
		zipVersion = fmt.Sprintf("%sdbru", version.Raw)
	}

	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "arm64",
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(
		urlPattern,
		c.downloadUrl,
		subFolder,
		archPart,
		zipVersion,
	), nil
}

func (c *instantClientComponent) GetAllVersions() ([]*gover.Version, error) {
	// Parse the latest versions from download page. The page may contain several links on a single
	// very long line, so matches are extracted from the whole content instead of on a per-line basis.
	content, err := installer.Tools.Download.AsBytes(c.versionsUrl)
	if err != nil {
		return nil, err
	}
	versions := []*gover.Version{}
	for _, match := range indexLineRegexp.FindAllStringSubmatch(string(content), -1) {
		versions = append(versions, gover.MustParseVersionFromRegex(match[2], versionRegexp))
	}
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
	defer os.Remove(fileName)
	rootFolder, err := installer.Tools.Compression.GetRootFolderFromZip(fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Compression.ExtractZip(fileName, "/opt/oracle", false); err != nil {
		return err
	}

	osInfo := installer.Tools.System.OsInfo()

	debianTrixieOrNewer := osInfo.IsDebian() && osInfo.MajorVersion() >= 13
	ubuntuNobleOrNewer := osInfo.IsUbuntu() && osInfo.MajorVersion() >= 24

	libaio1PackageName := "libaio1"
	if debianTrixieOrNewer || ubuntuNobleOrNewer {
		libaio1PackageName = "libaio1t64"
	}
	if err := installer.Tools.Apt.InstallDependencies(libaio1PackageName); err != nil {
		return err
	}

	if err := os.WriteFile("/etc/ld.so.conf.d/oracle-instantclient.conf", []byte(path.Join("/opt/oracle", rootFolder)), 0644); err != nil {
		return err
	}
	if err := goext.CmdRunners.Console.Run("ldconfig"); err != nil {
		return err
	}
	return nil
}
