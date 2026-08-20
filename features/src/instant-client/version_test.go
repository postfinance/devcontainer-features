package main

import (
	"builder/installer"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestInstantClientComponent() (*instantClientComponent, error) {
	versionsUrl, err := installer.Tools.System.MapArchitecture(map[string]string{
		installer.AMD64: defaultVersionsUrlAMD64,
		installer.ARM64: defaultVersionsUrlArm64,
	})
	if err != nil {
		return nil, err
	}

	return &instantClientComponent{
		ComponentBase: installer.NewComponentBase("Basic Package", "latest"),
		versionsUrl:   versionsUrl,
		downloadUrl:   defaultDownloadUrl,
	}, nil
}

func TestDownloadURLForAllMainVersions(t *testing.T) {
	c, err := newTestInstantClientComponent()
	assert.NoError(t, err)

	versions, err := c.GetAllVersions()
	assert.NoError(t, err)
	assert.NotEmpty(t, versions)

	for _, version := range versions {
		t.Run(version.String(), func(t *testing.T) {
			url, err := c.createDownloadURLForVersion(version)
			require.NoError(t, err)

			resp, err := http.Head(url)
			require.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, http.StatusOK, resp.StatusCode, "expected HTTP 200 for %s", url)
		})
	}
}
