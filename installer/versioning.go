package installer

import (
	"fmt"
	"regexp"

	"github.com/roemer/gover"
)

type versioning struct{}

func (v *versioning) ParseVersionsFromList(versions []string, versionRegex *regexp.Regexp, skipNonMatching bool) ([]*gover.Version, error) {
	parsedVersions := []*gover.Version{}
	for _, versionString := range versions {
		if versionRegex.MatchString(versionString) {
			version, err := gover.ParseVersionFromRegex(versionString, versionRegex)
			if err != nil {
				return nil, err
			}
			parsedVersions = append(parsedVersions, version)
		} else if !skipNonMatching {
			return nil, fmt.Errorf("version string '%s' does not match the expected format", versionString)
		}
	}
	return parsedVersions, nil
}
