package markup

import (
	"fmt"
	"path/filepath"

	"github.com/ci-space/edit-config/internal/fs"
)

func ParseDocument(file *fs.File) (Document, error) {
	ext := filepath.Ext(file.Path)
	switch ext {
	case ".yaml", ".yml":
		return LoadYAMLDocument(file.Content)
	default:
		return nil, fmt.Errorf("format %q not supported", ext)
	}
}
