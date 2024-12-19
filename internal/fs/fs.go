package fs

import (
	"bytes"
	"io"
)

type Filesystem interface {
	WriteFile(path string, content []byte) error
	ReadFile(path string) (*File, error)
}

type File struct {
	Path    string
	Content []byte
}

func (f *File) Reader() io.Reader {
	return bytes.NewReader(f.Content)
}
