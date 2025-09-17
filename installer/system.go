package installer

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
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

type OsInfo struct {
	Vendor    string
	Codename  string
	VersionId string
}

func (s *system) GetOsInfo() (*OsInfo, error) {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return nil, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	infoMap := map[string]string{}
	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), "=", 2)
		// Remove surrounding quotes if present
		val := strings.Trim(parts[1], `"`)
		infoMap[parts[0]] = val
	}
	return &OsInfo{
		Vendor:    infoMap["ID"],
		Codename:  infoMap["VERSION_CODENAME"],
		VersionId: infoMap["VERSION_ID"],
	}, nil
}
