package installer

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
)

type gitHub struct{}

func (g *gitHub) GetTags(owner string, repo string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/tags?per_page=100", owner, repo)
	nextRegexp := regexp.MustCompile(`(?i)<([^<]*)>; rel="next"`)
	var tagNames []string
	for {
		// Prepare the request
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, err
		}
		// Add the authorization header
		apiToken := os.Getenv("DEV_FEATURE_TOKEN_GITHUB_API")
		if apiToken != "" {
			req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiToken))
		}
		// Get the date for the current page
		resp, err := http.DefaultClient.Do(req)
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
