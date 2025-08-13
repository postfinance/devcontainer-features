package installer

type tools struct {
	Compression *compression
	Download    *download
	FileSystem  *fileSystem
	GitHub      *gitHub
	Http        *httpTools
}

var Tools *tools

func init() {
	Tools = &tools{
		Compression: &compression{},
		Download:    &download{},
		FileSystem:  &fileSystem{},
		GitHub:      &gitHub{},
		Http:        &httpTools{},
	}
}
