package actions

import (
	"fmt"
	"github.com/ci-space/edit-config/internal/fs"
)

type AppendAction struct {
	fs fs.Filesystem
}

func NewAppendAction(fs fs.Filesystem) *AppendAction {
	return &AppendAction{fs: fs}
}

func (act *AppendAction) Run(params Params) (*Result, error) {
	doc, err := loadDocument(act.fs, params.Filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load document: %v", err)
	}

	err = doc.Append(params.Pointer, params.NewValue)
	if err != nil {
		return nil, fmt.Errorf("failed to update values: %v", err)
	}

	newContent := doc.String()

	if params.DryRun {
		return act.dryRun(newContent, params)
	}

	err = act.fs.WriteFile(params.Filepath, []byte(newContent))
	if err != nil {
		return nil, fmt.Errorf("failed to write updated config: %w", err)
	}

	return &Result{
		Rows: []ResultRow{
			{
				Title: fmt.Sprintf("Updated file %s", params.Filepath),
			},
		},
	}, nil
}

func (act *AppendAction) dryRun(newContent string, params Params) (*Result, error) {
	return &Result{
		Rows: []ResultRow{
			{
				Title:   fmt.Sprintf("Dry Run for append %s", params.Pointer),
				Content: fmt.Sprintf("New content: \n%s", newContent),
			},
		},
	}, nil
}
