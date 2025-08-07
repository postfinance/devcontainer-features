package installer

type tools struct {
	Compression *compression
	Download    *download
	FileSystem  *fileSystem
}

var Tools *tools

func init() {
	Tools = &tools{
		Compression: &compression{},
		Download:    &download{},
		FileSystem:  &fileSystem{},
	}
}
