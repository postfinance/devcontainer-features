package installer

const (
	// The latest available version.
	VERSION_LATEST = "latest"
	// The latest lts version.
	VERSION_LTS = "lts"
	// No installation.
	VERSION_NONE = "none"
	// No installation as the included version from a previous component is good enough.
	VERSION_INCLUDED = "included"
	// Install the system default version.
	VERSION_SYSTEM_DEFAULT = "system-default"
	// The component does not have a version, so just install it.
	VERSION_IRRELEVANT = "irrelevant"
)
