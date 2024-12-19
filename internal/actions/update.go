package actions

import (
	"fmt"
	"github.com/ci-space/edit-config/internal/fs"

	"github.com/artarts36/yamlpath"
)

type UpdateAction struct {
	filesystem fs.Filesystem
}

func NewUpdateAction(fs fs.Filesystem) *UpdateAction {
	return &UpdateAction{filesystem: fs}
}

func (act *UpdateAction) Run(doc *yamlpath.Document, params Params) (*Result, error) {
	err := doc.Update(yamlpath.NewPointer(params.Pointer), params.NewValue)
	if err != nil {
		return nil, err
	}

	content, err := doc.Marshal()
	if err != nil {
		return nil, err
	}

	if params.DryRun {
		return act.dryRun(content, params)
	}

	err = act.filesystem.WriteFile(params.Filepath, content)
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

func (act *UpdateAction) dryRun(newContent []byte, params Params) (*Result, error) {
	return &Result{
		Rows: []ResultRow{
			{
				Title:   fmt.Sprintf("Dry Run for update %s", params.Pointer),
				Content: fmt.Sprintf("New content: \n%s", newContent),
			},
		},
	}, nil
}
