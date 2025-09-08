package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/roemer/gover"
)

type gitHub struct{}

func (g *gitHub) GetTags(owner string, repo string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?per_page=100", owner, repo)
	nextRegexp := regexp.MustCompile(`(?i)<([^<]*)>; rel="next"`)
	var tagNames []string
	for {
		// Get the date for the current page
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("failed to download file '%s'. Status code: %d", url, resp.StatusCode)
		}
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read body: %w", err)
		}
		// Parse the items
		var pageItems []*struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal(bodyBytes, &pageItems); err != nil {
			return nil, err
		}
		// Add the items to the result list
		for _, item := range pageItems {
			tagNames = append(tagNames, item.Name)
		}
		// Search for a next link
		if linkHeader, ok := resp.Header["Link"]; ok {
			matches := nextRegexp.FindStringSubmatch(linkHeader[0])
			if matches != nil {
				// Set the new url and continue the loop
				url = matches[1]
				continue
			}
		}
		// No next link, abort the loop
		break
	}
	// Return the found items
	return tagNames, nil
}

func (g *gitHub) ParseVersionFromTags(tags []string, versionRegex *regexp.Regexp, excludeTags ...string) ([]*gover.Version, error) {
	var versions []*gover.Version
	for _, tag := range tags {
		// Skip tags in excludeTags
		if len(excludeTags) > 0 {
			for _, exclude := range excludeTags {
				if tag == exclude {
					continue
				}
			}
		}
		// Find version using the provided regex
		matches := versionRegex.FindStringSubmatch(tag)
		if len(matches) > 0 {
			verStr := matches[0]
			ver := gover.MustParseVersionFromRegex(verStr, versionRegex)
			versions = append(versions, ver)
		}
	}
	if len(versions) == 0 {
		return nil, fmt.Errorf("no valid versions found in tags")
	}
	return versions, nil
}
