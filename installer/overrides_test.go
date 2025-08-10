package installer

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOverrideUseDefault(t *testing.T) {
	assert := assert.New(t)

	var testValue string
	HandleOverride(&testValue, "default", "test-key")
	assert.Equal("default", testValue)
}

func TestOverrideUseEnv(t *testing.T) {
	assert := assert.New(t)

	var testValue string
	os.Setenv("DEV_FEATURE_OVERRIDE_TEST_KEY", "override")
	HandleOverride(&testValue, "default", "test-key")
	assert.Equal("override", testValue)
}

func TestOverrideUseFile(t *testing.T) {
	assert := assert.New(t)

	// Setup http server for an override file
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)
	mux.HandleFunc("/file", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("DEV_FEATURE_OVERRIDE_TEST_KEY=overridefile\nTEST_KEY_2=overridefile2\nDEV_FEATURE_OVERRIDE_TEST_KEY_SET=yes-override"))
	})

	// Manually set an env value beforehand
	os.Setenv("DEV_FEATURE_OVERRIDE_TEST_KEY_SET", "no-override")

	// Load the overrides of the file
	os.Setenv("DEV_FEATURE_OVERRIDE_LOCATION", server.URL+"/file")
	loadErr := LoadOverrides()
	assert.NoError(loadErr, "LoadOverrides should not return an error")

	// Test the values
	var testValue string
	HandleOverride(&testValue, "default", "test-key")
	assert.Equal("overridefile", testValue)

	var testValue2 string
	HandleOverride(&testValue2, "default", "test-key-2")
	assert.Equal("overridefile2", testValue2)

	var testValue3 string
	HandleOverride(&testValue3, "default", "test-key-set")
	assert.Equal("no-override", testValue3)
}
