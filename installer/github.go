package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
)

type gitHub struct{}

func (g *gitHub) GetTags(owner string, repo string) ([]*GitHubTag, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?per_page=100", owner, repo)
	nextRegexp := regexp.MustCompile(`(?i)<([^<]*)>; rel="next"`)

	var gitHubTags []*GitHubTag
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
		var pageItems []*GitHubTag
		if err := json.Unmarshal(bodyBytes, &pageItems); err != nil {
			return nil, err
		}
		gitHubTags = append(gitHubTags, pageItems...)
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
	return gitHubTags, nil
}

type GitHubTag struct {
	Name string `json:"name"`
}
