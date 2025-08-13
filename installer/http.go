package installer

import (
	"bufio"
	"bytes"
	"fmt"
	"net/url"
	"regexp"

	"github.com/roemer/gover"
)

type httpTools struct{}

func (h *httpTools) BuildUrl(base string, parts ...string) (string, error) {
	return url.JoinPath(base, parts...)
}

// Function based variant to get versions from a html page with custom parsing.
func (h *httpTools) GetVersionsFromHtmlIndexFunc(url string, lineFunc func(url string, line string) (*gover.Version, error)) ([]*gover.Version, error) {
	versionFileContent, err := Tools.Download.AsBytes(url)
	if err != nil {
		return nil, err
	}
	allVersions := []*gover.Version{}
	scanner := bufio.NewScanner(bytes.NewReader(versionFileContent))
	for scanner.Scan() {
		version, err := lineFunc(url, scanner.Text())
		if err != nil {
			return nil, err
		}
		if version != nil {
			allVersions = append(allVersions, version)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed processing the line scanner")
	}
	return allVersions, err
}

// Simple variant to get versions from the lines of a html page.
func (h *httpTools) GetVersionsFromHtmlIndex(url string, lineRegex *regexp.Regexp, versionRegex *regexp.Regexp) ([]*gover.Version, error) {
	return h.GetVersionsFromHtmlIndexFunc(url, func(url, line string) (*gover.Version, error) {
		if match := lineRegex.FindStringSubmatch(line); match != nil {
			versionString := match[1]
			return gover.ParseVersionFromRegex(versionString, versionRegex)
		}
		return nil, nil
	})
}
