package installer

import (
	"fmt"
	"runtime"
)

const (
	AMD64 = "amd64"
	ARM64 = "arm64"
)

type system struct{}

func (s *system) MapArchitecture(mapping map[string]string) (string, error) {
	mappedValue, ok := mapping[runtime.GOARCH]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}
	return mappedValue, nil
}
