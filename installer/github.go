package installer

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/roemer/goext"
)

const headerRateLimitRemaining = "x-ratelimit-remaining"

type gitHub struct{}

func (g *gitHub) GetTags(owner string, repo string) ([]string, error) {
	// fetching tags via GitHub API...
	tags, err := g.fetchTagsViaGithubAPI(owner, repo)
	if err == nil {
		// ...was successful -> return the tags
		return tags, nil
	}

	// ...failed
	var rateLimitErr *GitHubRateLimitError
	// check if err is a GitHubRateLimitError
	if ok := errors.As(err, &rateLimitErr); ok {
		// err is a GitHubRateLimitError
		// should we use the fallback?
		if rateLimitErr.UseFallback {
			// yes
			return g.fetchTagsViaGitLsRemote(owner, repo)
		} else {
			// no, return the error
			return nil, err
		}
	}

	// if the error is not a GitHubRateLimitError, we will return the error
	return nil, err
}

type GitHubRateLimitError struct {
	UseFallback bool
	Message     string
}

func (e *GitHubRateLimitError) Error() string {
	return fmt.Sprintf("GitHub API rate limit error: %s", e.Message)
}

func (g *gitHub) fetchTagsViaGithubAPI(owner string, repo string) ([]string, error) {
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
			return g.handleRateLimitError(resp, url, apiToken)
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

func (g *gitHub) fetchTagsViaGitLsRemote(owner, repo string) ([]string, error) {
	gitTagsStdout, gitTagsStderr, err := goext.CmdRunners.Default.RunGetOutput("git", "ls-remote", "--tags", fmt.Sprintf("https://github.com/%s/%s.git", owner, repo))
	if err != nil {
		return nil, fmt.Errorf("failed getting git tags: %v - %s", err, gitTagsStderr)
	}

	lineRegexp := regexp.MustCompile(`^[a-fA-F0-9]+\s+refs/tags/(.*)$`)
	var tagNames []string
	scanner := bufio.NewScanner(strings.NewReader(gitTagsStdout))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasSuffix(line, "^{}") {
			continue
		}
		if matches := lineRegexp.FindStringSubmatch(line); matches != nil {
			tagName := matches[1]
			tagNames = append(tagNames, tagName)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("failed processing the line scanner: %w", err)
	}

	return tagNames, nil
}

func (g *gitHub) handleRateLimitError(resp *http.Response, url, apiToken string) ([]string, error) {
	if !g.isRateLimit(resp) {
		return nil, fmt.Errorf("failed to download file '%s'. Status code: %d", url, resp.StatusCode)
	}

	hasToken := apiToken != ""
	if !hasToken {
		return nil, &GitHubRateLimitError{
			Message: fmt.Sprintf(
				"failed to download file '%s'. Status code: %d. You may have reached the GitHub API rate limit. Consider setting the DEV_FEATURE_TOKEN_GITHUB_API environment variable with a personal access token to increase the limit.",
				url, resp.StatusCode,
			),
			UseFallback: false,
		}
	}

	return nil, &GitHubRateLimitError{
		Message: fmt.Sprintf(
			"failed to download file '%s'. Status code: %d. You may have reached the GitHub API rate limit even with a personal access token. A fallback via git ls-remote is attempted...",
			url, resp.StatusCode,
		),
		UseFallback: true,
	}
}

func (g *gitHub) isRateLimit(resp *http.Response) bool {
	// rate limit rules: https://docs.github.com/en/rest/using-the-rest-api/rate-limits-for-the-rest-api?apiVersion=2026-03-10
	// the method net/http.Header.Get() already canonicalize the header name, so we can use it directly
	if (resp.StatusCode == http.StatusForbidden || resp.StatusCode == http.StatusTooManyRequests) && resp.Header.Get(headerRateLimitRemaining) == "0" {
		return true
	}
	return false
}
