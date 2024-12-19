package fs

type Filesystem interface {
	WriteFile(path string, content []byte) error
}
