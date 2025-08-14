package installer

import (
	"fmt"
	"strings"

	"github.com/roemer/gover"
)

// Constructor for a feature.
func NewFeature(name string, ensureLoginShellPath bool, components ...IComponent) *Feature {
	return &Feature{
		name:                 name,
		ensureLoginShellPath: ensureLoginShellPath,
		components:           components,
	}
}

// This type represents a feature that contains one or more components.
type Feature struct {
	name                 string
	ensureLoginShellPath bool
	components           []IComponent
}

// Adds one or more components to the feature.
func (f *Feature) AddComponents(components ...IComponent) {
	f.components = append(f.components, components...)
}

// Main function that processes a given feature.
func (f *Feature) Process() error {
	fmt.Printf("Processing feature '%s'\n", f.name)

	// Process the components
	for _, component := range f.components {
		fmt.Printf("Processing component '%s'\n", component.GetName())
		requestedVersionString := component.GetRequestedVersion()
		isExactVersion := strings.HasSuffix(requestedVersionString, "!")
		requestedVersionString = strings.TrimSuffix(requestedVersionString, "!")

		// Skip if "none" version was requested
		if requestedVersionString == VERSION_NONE {
			fmt.Println("Version is none, skipping")
			continue
		}

		// Skip if "included" version was requested
		if requestedVersionString == VERSION_INCLUDED {
			fmt.Println("Version is included, skipping")
			continue
		}

		// Directly install it if the version is system default or irrelevant
		if requestedVersionString == VERSION_SYSTEM_DEFAULT || requestedVersionString == VERSION_IRRELEVANT {
			fmt.Println("Installing default")
			if err := component.InstallVersion(gover.EmptyVersion); err != nil {
				return err
			}
			continue
		}

		// Calculate the version to install
		var versionToInstall *gover.Version
		// Handle non-numeric (latest/lts)
		if requestedVersionString == VERSION_LATEST || requestedVersionString == VERSION_LTS {
			// Try get the latest version
			version, err := component.GetLatestVersion()
			if err != nil {
				return err
			}
			if version != nil {
				// Found it and use it
				versionToInstall = version
			} else {
				// Default to the max of all versions
				allVersions, err := component.GetAllVersions()
				if err != nil {
					return err
				}
				versionToInstall = gover.FindMax(allVersions, gover.EmptyVersion, true)
			}
		} else {
			// Parse the reference version
			referenceVersion := gover.ParseSimple(strings.Split(strings.ReplaceAll(requestedVersionString, "-", "."), "."))
			if isExactVersion {
				// The exact version was passed, so directly use it
				versionToInstall = referenceVersion
				versionToInstall.Raw = requestedVersionString
			} else {
				// Get all versions
				allVersions, err := component.GetAllVersions()
				if err != nil {
					return err
				}
				// Get the max according to the reference version
				versionToInstall = gover.FindMax(allVersions, referenceVersion, false)
			}
		}

		// No version found
		if versionToInstall == nil {
			return fmt.Errorf("no version to install found for '%s'", requestedVersionString)
		}

		// Version found, install it
		fmt.Printf("Installing version %s\n", versionToInstall.Raw)
		if err := component.InstallVersion(versionToInstall); err != nil {
			return err
		}
	}

	if f.ensureLoginShellPath {
		// Ensure that login shells get the correct path if the user updated the PATH using ENV.
		if err := Tools.FileSystem.EnsureLoginShellPath(); err != nil {
			return err
		}
	}

	return nil
}

// This is the inferface that needs to be implemented by a component.
type IComponent interface {
	// Returns the name of the component.
	GetName() string
	// Gets the requestedVersion of the component.
	GetRequestedVersion() string
	// Returns a list of all available versions.
	GetAllVersions() ([]*gover.Version, error)
	// Returns the latest version. Defaults to nil which then uses the max of GetAllVersions.
	GetLatestVersion() (*gover.Version, error)
	// Installs the given version.
	InstallVersion(version *gover.Version) error
}

// Constructor for a base component.
func NewComponentBase(name string, requestedVersion string) *ComponentBase {
	return &ComponentBase{
		name:             name,
		requestedVersion: requestedVersion,
	}
}

// Base struct for a component with common things.
type ComponentBase struct {
	name             string
	requestedVersion string
}

// Gets the name of the component.
func (c *ComponentBase) GetName() string {
	return c.name
}

// Gets the requestedVersion of the component.
func (c *ComponentBase) GetRequestedVersion() string {
	return c.requestedVersion
}

// Gets all possible version. Returns nil if no implementation is provided.
func (c *ComponentBase) GetAllVersions() ([]*gover.Version, error) {
	return nil, nil
}

// Gets the latest version. Returns nil if no implementation is provided.
func (c *ComponentBase) GetLatestVersion() (*gover.Version, error) {
	return nil, nil
}
