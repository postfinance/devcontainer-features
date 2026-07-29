package installer

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myFakeService func(req *http.Request) (*http.Response, error)

func (f myFakeService) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func mockDefaultHTTPClient(t *testing.T, rt myFakeService) {
	t.Helper()
	oldClient := http.DefaultClient
	http.DefaultClient = &http.Client{Transport: rt}
	t.Cleanup(func() {
		http.DefaultClient = oldClient
	})
}

func TestFetchTagsViaGithubAPIWithoutTokenSuccess(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("DEV_FEATURE_TOKEN_GITHUB_API", "")

	var authorizationHeader string
	mockDefaultHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		authorizationHeader = req.Header.Get("Authorization")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"name":"v1.0.0"},{"name":"v1.1.0"}]`)),
		}, nil
	})

	g := &gitHub{}
	tags, err := g.fetchTagsViaGithubAPI("owner", "repo")
	require.NoError(t, err)
	assert.Empty(t, authorizationHeader)
	assert.Contains(t, tags, "v1.1.0")
}

func TestFetchTagsViaGithubAPIWithTokenSuccess(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("DEV_FEATURE_TOKEN_GITHUB_API", "ghp_dummy_token")

	var authorizationHeader string
	mockDefaultHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		authorizationHeader = req.Header.Get("Authorization")
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     make(http.Header),
			Body:       io.NopCloser(strings.NewReader(`[{"name":"v2.0.0"},{"name":"v2.1.0"}]`)),
		}, nil
	})

	g := &gitHub{}
	tags, err := g.fetchTagsViaGithubAPI("owner", "repo")
	require.NoError(t, err)
	assert.Equal(t, "Bearer ghp_dummy_token", authorizationHeader)
	assert.Contains(t, tags, "v2.1.0")
}

func TestFetchTagsViaGitLsRemoteSuccess(t *testing.T) {
	t.Chdir(t.TempDir())

	tmpDir := t.TempDir()
	fakeGitPath := filepath.Join(tmpDir, "git")
	err := os.WriteFile(fakeGitPath, []byte(`#!/usr/bin/env sh
if [ "$1" = "ls-remote" ] && [ "$2" = "--tags" ]; then
  cat <<'EOF'
1111111111111111111111111111111111111111	refs/tags/v0.1.0
1111111111111111111111111111111111111111	refs/tags/v0.1.0^{}
2222222222222222222222222222222222222222	refs/tags/v0.2.0
EOF
  exit 0
fi
echo "unexpected args: $@" >&2
exit 1
`), 0755)
	require.NoError(t, err)

	path := os.Getenv("PATH")
	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+path)

	g := &gitHub{}
	tags, err := g.fetchTagsViaGitLsRemote("owner", "repo")
	require.NoError(t, err)
	assert.Contains(t, tags, "v0.1.0")
	assert.Contains(t, tags, "v0.2.0")
	assert.Len(t, tags, 2)
}

func TestFetchTagsViaGithubAPIWithoutTokenFailed(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("DEV_FEATURE_TOKEN_GITHUB_API", "")

	mockDefaultHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		h := make(http.Header)
		h.Set(headerRateLimitRemaining, "0")
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     h,
			Body:       io.NopCloser(strings.NewReader(`{"message":"API rate limit exceeded"}`)),
		}, nil
	})

	g := &gitHub{}
	tags, err := g.fetchTagsViaGithubAPI("owner", "repo")
	require.Error(t, err)
	assert.Nil(t, tags)

	var rateLimitErr *GitHubRateLimitError
	require.ErrorAs(t, err, &rateLimitErr)
	assert.False(t, rateLimitErr.UseFallback)
	assert.Contains(t, err.Error(), "You may have reached the GitHub API rate limit")
	assert.Contains(t, err.Error(), "DEV_FEATURE_TOKEN_GITHUB_API")
}

func TestFetchTagsViaGithubAPIWithTokenFailed(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("DEV_FEATURE_TOKEN_GITHUB_API", "ghp_dummy_token")

	mockDefaultHTTPClient(t, func(req *http.Request) (*http.Response, error) {
		h := make(http.Header)
		h.Set(headerRateLimitRemaining, "0")
		return &http.Response{
			StatusCode: http.StatusForbidden,
			Header:     h,
			Body:       io.NopCloser(strings.NewReader(`{"message":"API rate limit exceeded"}`)),
		}, nil
	})

	g := &gitHub{}
	tags, err := g.fetchTagsViaGithubAPI("owner", "repo")
	require.Error(t, err)
	assert.Nil(t, tags)

	var rateLimitErr *GitHubRateLimitError
	require.ErrorAs(t, err, &rateLimitErr)
	assert.True(t, rateLimitErr.UseFallback)
	assert.Contains(t, err.Error(), "You may have reached the GitHub API rate limit even with a personal access token")
	assert.Contains(t, err.Error(), "A fallback via git ls-remote is attempted")
}
