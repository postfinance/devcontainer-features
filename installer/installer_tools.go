package installer

type tools struct {
	Compression *compression
	Download    *download
	FileSystem  *fileSystem
	GitHub      *gitHub
	GitLab      *gitLab
	Http        *httpTools
	Npm         *npm
	System      *system
	Versioning  *versioning
	Apt         *apt
}

var Tools *tools

func init() {
	Tools = &tools{
		Compression: &compression{},
		Download:    &download{},
		FileSystem:  &fileSystem{},
		GitHub:      &gitHub{},
		GitLab:      &gitLab{},
		Http:        &httpTools{},
		Npm:         &npm{},
		System:      &system{},
		Versioning:  &versioning{},
		Apt:         &apt{},
	}
}
