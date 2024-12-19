package fs

import (
	"os"
)

type Local struct{}

func NewLocal() *Local {
	return &Local{}
}

func (l *Local) WriteFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}

func (l *Local) ReadFile(path string) (*File, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return &File{
		Path:    path,
		Content: content,
	}, nil
}
