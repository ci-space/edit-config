package fs

import "os"

type Local struct{}

func NewLocal() *Local {
	return &Local{}
}

func (l *Local) WriteFile(path string, content []byte) error {
	return os.WriteFile(path, content, 0644)
}
