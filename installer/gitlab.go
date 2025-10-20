package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

type gitLab struct{}

func (g *gitLab) GetPackageReleases(projectPath string) ([]string, error) {
	encodedProjectPath := url.PathEscape(projectPath)
	url := fmt.Sprintf("https://gitlab.com/api/v4/projects/%s/packages?order_by=created_at&sort=desc&per_page=100", encodedProjectPath)
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
			Id        int       `json:"id"`
			Name      string    `json:"name"`
			Version   string    `json:"version"`
			CreatedAt time.Time `json:"created_at"`
		}
		if err := json.Unmarshal(bodyBytes, &pageItems); err != nil {
			return nil, err
		}
		// Add the items to the result list
		for _, item := range pageItems {
			tagNames = append(tagNames, item.Version)
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
