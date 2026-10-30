package main

import (
	"testing"

	"github.com/roemer/gover"
)

func TestParseNodeVersion(t *testing.T) {
	if _, err := gover.ParseVersionFromRegex("v27.0.0-alpha.0.0.0", nodeVersionRegex); err == nil {
		t.Fatal("expected prerelease version not to match")
	}

	version, err := gover.ParseVersionFromRegex("v27.0.0", nodeVersionRegex)
	if err != nil {
		t.Fatal("expected stable version to match")
	}
	if version.Raw != "27.0.0" {
		t.Fatalf("expected raw version 27.0.0, got %s", version.Raw)
	}
}
