package actions

import (
	"fmt"

	"github.com/ci-space/edit-config/internal/fs"
	"github.com/ci-space/edit-config/internal/shared/markup"
)

func loadDocument(fs fs.Filesystem, filepath string) (markup.Document, error) {
	file, err := fs.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	doc, err := markup.ParseDocument(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse document: %v", err)
	}

	return doc, nil
}
