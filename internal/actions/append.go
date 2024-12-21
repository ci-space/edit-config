package actions

import (
	"fmt"
	"strings"

	"github.com/ci-space/edit-config/internal/fs"
	"github.com/ci-space/edit-config/internal/shared/markup"
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

	appended, err := act.append(params, doc)
	if err != nil {
		return nil, fmt.Errorf("failed to append document: %v", err)
	}

	newContent := doc.String()

	if params.DryRun {
		return act.dryRun(newContent, params, appended)
	}

	err = act.fs.WriteFile(params.Filepath, []byte(newContent))
	if err != nil {
		return nil, fmt.Errorf("failed to write updated config: %w", err)
	}

	return &Result{
		Rows: []ResultRow{
			{
				Title: fmt.Sprintf("Updated file %s, appended %d items", params.Filepath, appended),
			},
		},
	}, nil
}

func (act *AppendAction) append(params Params, doc markup.Document) (int, error) {
	var err error

	if params.Separator == "" {
		err = doc.Append(params.Pointer, params.NewValue)
		if err != nil {
			return 0, fmt.Errorf("failed to update values: %v", err)
		}
		return 1, nil
	}

	splitted := strings.Split(params.NewValue, params.Separator)

	for _, v := range splitted {
		err = doc.Append(params.Pointer, v)
		if err != nil {
			return 0, fmt.Errorf("failed to update values: %v", err)
		}
	}

	return len(splitted), nil
}

func (act *AppendAction) dryRun(newContent string, params Params, appended int) (*Result, error) {
	return &Result{
		Rows: []ResultRow{
			{
				Title:   fmt.Sprintf("Dry Run: file %s updated, appended %d items", params.Pointer, appended),
				Content: fmt.Sprintf("New content: \n%s", newContent),
			},
		},
	}, nil
}
