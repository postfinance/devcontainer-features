package main

import (
	"builder/installer"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"regexp"

	"github.com/roemer/gover"
)

//////////
// Configuration
//////////

// Version regex for Temurin/OpenJDK releases like "jdk-21.0.7+6"
var javaVersionRegex *regexp.Regexp = regexp.MustCompile(`^jdk-(?P<raw>(\d+)\.(\d+)\.(\d+))(?:\+\d+)?$`)

// Version regex for Maven releases like "3.9.9"
var mavenVersionRegex *regexp.Regexp = regexp.MustCompile(`^(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

// Version regex for Gradle releases like "8.14" or "8.2.1"
var gradleVersionRegex *regexp.Regexp = regexp.MustCompile(`^(?P<raw>(\d+)\.(\d+)(?:\.(\d+))?)$`)

// Version regex for Ant releases like "1.10.15"
var antVersionRegex *regexp.Regexp = regexp.MustCompile(`^(?P<raw>(\d+)\.(\d+)\.(\d+))$`)

// Line regex for parsing Maven HTML index pages like href="3.9.9/"
var mavenIndexLineRegex *regexp.Regexp = regexp.MustCompile(`href="(\d+\.\d+\.\d+)/"`)

// Line regex for parsing Ant HTML index pages like href="apache-ant-1.10.15-bin.tar.gz"
var antIndexLineRegex *regexp.Regexp = regexp.MustCompile(`href="apache-ant-(\d+\.\d+\.\d+)-bin\.tar\.gz"`)

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
	mavenVersion := flag.String("mavenVersion", "none", "")
	gradleVersion := flag.String("gradleVersion", "none", "")
	antVersion := flag.String("antVersion", "none", "")
	downloadUrl := flag.String("downloadUrl", "", "")
	versionsUrl := flag.String("versionsUrl", "", "")
	latestUrl := flag.String("latestUrl", "", "")
	mavenDownloadUrl := flag.String("mavenDownloadUrl", "", "")
	mavenVersionsUrl := flag.String("mavenVersionsUrl", "", "")
	gradleDownloadUrl := flag.String("gradleDownloadUrl", "", "")
	gradleVersionsUrl := flag.String("gradleVersionsUrl", "", "")
	antDownloadUrl := flag.String("antDownloadUrl", "", "")
	antVersionsUrl := flag.String("antVersionsUrl", "", "")
	flag.Parse()

	// Load settings from an external file
	if err := installer.LoadOverrides(); err != nil {
		return err
	}

	installer.HandleOverride(downloadUrl, "https://api.adoptium.net/v3/binary", "java-download-url")
	installer.HandleOverride(versionsUrl, "https://api.adoptium.net/v3/info/release_names", "java-versions-url")
	installer.HandleOverride(latestUrl, "https://api.adoptium.net/v3/info/available_releases", "java-latest-url")
	installer.HandleOverride(mavenDownloadUrl, "https://downloads.apache.org/maven/maven-3", "java-maven-download-url")
	installer.HandleOverride(mavenVersionsUrl, "https://downloads.apache.org/maven/maven-3/", "java-maven-versions-url")
	installer.HandleOverride(gradleDownloadUrl, "https://services.gradle.org/distributions", "java-gradle-download-url")
	installer.HandleOverride(gradleVersionsUrl, "https://services.gradle.org/versions/all", "java-gradle-versions-url")
	installer.HandleOverride(antDownloadUrl, "https://downloads.apache.org/ant/binaries", "java-ant-download-url")
	installer.HandleOverride(antVersionsUrl, "https://downloads.apache.org/ant/binaries/", "java-ant-versions-url")

	// Create and process the feature
	feature := installer.NewFeature("Java", true,
		&javaComponent{
			ComponentBase: installer.NewComponentBase("Java", *version),
			DownloadUrl:   *downloadUrl,
			VersionsUrl:   *versionsUrl,
			LatestUrl:     *latestUrl,
			releaseNames:  make(map[string]string),
		},
		&mavenComponent{
			ComponentBase: installer.NewComponentBase("Maven", *mavenVersion),
			DownloadUrl:   *mavenDownloadUrl,
			VersionsUrl:   *mavenVersionsUrl,
		},
		&gradleComponent{
			ComponentBase: installer.NewComponentBase("Gradle", *gradleVersion),
			DownloadUrl:   *gradleDownloadUrl,
			VersionsUrl:   *gradleVersionsUrl,
		},
		&antComponent{
			ComponentBase: installer.NewComponentBase("Ant", *antVersion),
			DownloadUrl:   *antDownloadUrl,
			VersionsUrl:   *antVersionsUrl,
		},
	)
	return feature.Process()
}

//////////
// Java Component
//////////

type javaComponent struct {
	*installer.ComponentBase
	DownloadUrl  string
	VersionsUrl  string
	LatestUrl    string
	releaseNames map[string]string
}

// Always resolve through GetAllVersions to populate the release name map.
func (c *javaComponent) IsFullVersion(referenceVersion *gover.Version) bool {
	return false
}

func (c *javaComponent) GetLatestVersion() (*gover.Version, error) {
	// Get available releases info to determine the most recent LTS major version
	data, err := installer.Tools.Download.AsBytes(c.LatestUrl)
	if err != nil {
		return nil, err
	}
	var releasesInfo struct {
		MostRecentLts int `json:"most_recent_lts"`
	}
	if err := json.Unmarshal(data, &releasesInfo); err != nil {
		return nil, err
	}
	if releasesInfo.MostRecentLts == 0 {
		return nil, nil
	}

	// Get the latest release for the most recent LTS major version
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return nil, err
	}
	pageUrl := fmt.Sprintf("%s?release_type=ga&os=linux&architecture=%s&image_type=jdk&vendor=eclipse&feature_version=%d&page=0&page_size=1",
		c.VersionsUrl, archPart, releasesInfo.MostRecentLts)
	data, err = installer.Tools.Download.AsBytes(pageUrl)
	if err != nil {
		return nil, err
	}
	var response struct {
		Releases []string `json:"releases"`
	}
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, err
	}
	if len(response.Releases) == 0 {
		return nil, nil
	}
	version, err := gover.ParseVersionFromRegex(response.Releases[0], javaVersionRegex)
	if err != nil {
		return nil, err
	}
	c.releaseNames[version.Raw] = response.Releases[0]
	return version, nil
}

func (c *javaComponent) GetAllVersions() ([]*gover.Version, error) {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return nil, err
	}

	allVersions := []*gover.Version{}
	pageSize := 100
	for page := 0; ; page++ {
		pageUrl := fmt.Sprintf("%s?release_type=ga&os=linux&architecture=%s&image_type=jdk&vendor=eclipse&page=%d&page_size=%d",
			c.VersionsUrl, archPart, page, pageSize)
		data, err := installer.Tools.Download.AsBytes(pageUrl)
		if err != nil {
			return nil, err
		}
		var response struct {
			Releases []string `json:"releases"`
		}
		if err := json.Unmarshal(data, &response); err != nil {
			return nil, err
		}
		if len(response.Releases) == 0 {
			break
		}
		for _, releaseName := range response.Releases {
			version, err := gover.ParseVersionFromRegex(releaseName, javaVersionRegex)
			if err != nil {
				continue // skip e.g. Java 8 releases with a different naming scheme
			}
			c.releaseNames[version.Raw] = releaseName
			allVersions = append(allVersions, version)
		}
		if len(response.Releases) < pageSize {
			break // last page reached
		}
	}
	return allVersions, nil
}

func (c *javaComponent) InstallVersion(version *gover.Version) error {
	archPart, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: "x64",
		installer.ARM64: "aarch64",
	})
	if err != nil {
		return err
	}

	releaseName, ok := c.releaseNames[version.Raw]
	if !ok {
		return fmt.Errorf("could not find release name for Java version %s", version.Raw)
	}

	// Build the download URL (url.JoinPath encodes '+' as '%2B' automatically)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl,
		"version", releaseName, "linux", archPart, "jdk", "hotspot", "normal", "eclipse")
	if err != nil {
		return err
	}

	fileName := fmt.Sprintf("java-%s-linux-%s.tar.gz", version.Raw, archPart)
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Java"); err != nil {
		return err
	}

	// Remove old installation and extract to /usr/local/java (strip the root folder)
	os.RemoveAll("/usr/local/java")
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "/usr/local/java", true); err != nil {
		return err
	}

	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}

//////////
// Maven Component
//////////

type mavenComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
}

func (c *mavenComponent) GetAllVersions() ([]*gover.Version, error) {
	return installer.Tools.Http.GetVersionsFromHtmlIndex(c.VersionsUrl, mavenIndexLineRegex, mavenVersionRegex)
}

func (c *mavenComponent) InstallVersion(version *gover.Version) error {
	fileName := fmt.Sprintf("apache-maven-%s-bin.tar.gz", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, version.Raw, "binaries", fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Maven"); err != nil {
		return err
	}

	// Remove old installation and extract to /usr/local/maven (strip the root folder)
	os.RemoveAll("/usr/local/maven")
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "/usr/local/maven", true); err != nil {
		return err
	}

	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}

//////////
// Gradle Component
//////////

type gradleVersionEntry struct {
	Version        string `json:"version"`
	Nightly        bool   `json:"nightly"`
	Snapshot       bool   `json:"snapshot"`
	ReleaseNightly bool   `json:"releaseNightly"`
	ActiveRc       bool   `json:"activeRc"`
	Broken         bool   `json:"broken"`
	RcFor          string `json:"rcFor"`
	MilestoneFor   string `json:"milestoneFor"`
}

type gradleComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
}

// Always resolve through GetAllVersions to handle mixed 2/3-segment Gradle versions.
func (c *gradleComponent) IsFullVersion(referenceVersion *gover.Version) bool {
	return false
}

func (c *gradleComponent) GetAllVersions() ([]*gover.Version, error) {
	data, err := installer.Tools.Download.AsBytes(c.VersionsUrl)
	if err != nil {
		return nil, err
	}
	var entries []gradleVersionEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	for _, entry := range entries {
		// Only include stable (non-broken, non-pre-release) versions
		if entry.Broken || entry.Nightly || entry.Snapshot || entry.ReleaseNightly || entry.ActiveRc || entry.RcFor != "" || entry.MilestoneFor != "" {
			continue
		}
		version, err := gover.ParseVersionFromRegex(entry.Version, gradleVersionRegex)
		if err != nil {
			continue
		}
		allVersions = append(allVersions, version)
	}
	return allVersions, nil
}

func (c *gradleComponent) InstallVersion(version *gover.Version) error {
	fileName := fmt.Sprintf("gradle-%s-bin.zip", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Gradle"); err != nil {
		return err
	}

	// Remove old installation and extract to /usr/local/gradle (strip the root folder)
	os.RemoveAll("/usr/local/gradle")
	if err := installer.Tools.Compression.ExtractZip(fileName, "/usr/local/gradle", true); err != nil {
		return err
	}

	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}

//////////
// Ant Component
//////////

type antComponent struct {
	*installer.ComponentBase
	DownloadUrl string
	VersionsUrl string
}

func (c *antComponent) GetAllVersions() ([]*gover.Version, error) {
	return installer.Tools.Http.GetVersionsFromHtmlIndex(c.VersionsUrl, antIndexLineRegex, antVersionRegex)
}

func (c *antComponent) InstallVersion(version *gover.Version) error {
	fileName := fmt.Sprintf("apache-ant-%s-bin.tar.gz", version.Raw)
	downloadUrl, err := installer.Tools.Http.BuildUrl(c.DownloadUrl, fileName)
	if err != nil {
		return err
	}
	if err := installer.Tools.Download.ToFile(downloadUrl, fileName, "Ant"); err != nil {
		return err
	}

	// Remove old installation and extract to /usr/local/ant (strip the root folder)
	os.RemoveAll("/usr/local/ant")
	if err := installer.Tools.Compression.ExtractTarGz(fileName, "/usr/local/ant", true); err != nil {
		return err
	}

	// Cleanup
	if err := os.Remove(fileName); err != nil {
		return err
	}
	return nil
}
