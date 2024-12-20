package actions

import (
	"fmt"
	"strings"

	githuboutput "github.com/ci-space/github-output"

	"github.com/ci-space/edit-config/internal/fs"
	"github.com/ci-space/edit-config/internal/shared/image"
)

type UpImageVersionAction struct {
	fs fs.Filesystem
}

func NewUpImageVersionAction(fs fs.Filesystem) *UpImageVersionAction {
	return &UpImageVersionAction{fs: fs}
}

func (act *UpImageVersionAction) Run(params Params) (*Result, error) {
	doc, err := loadDocument(act.fs, params.Filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to load document: %v", err)
	}

	var imageString string

	err = doc.Read(params.Pointer, &imageString)
	if err != nil {
		return nil, fmt.Errorf("failed to read current version: %v", err)
	}

	vImage, err := image.ParseImage(imageString)
	if err != nil {
		return nil, err
	}

	err = act.upVersion(vImage, params)
	if err != nil {
		return nil, err
	}

	err = doc.UpdateValue(params.Pointer, strings.NewReader(vImage.String()))
	if err != nil {
		return nil, fmt.Errorf("failed to update version: %v", err)
	}

	newContent := doc.String()

	if params.DryRun {
		return act.dryRun(newContent, params)
	}

	err = act.fs.WriteFile(params.Filepath, []byte(newContent))
	if err != nil {
		return nil, fmt.Errorf("failed to write updated config: %w", err)
	}

	err = githuboutput.WhenAvailable(func() error {
		return githuboutput.Write("new-version", vImage.Version.String())
	})
	if err != nil {
		return nil, fmt.Errorf("failed to write to github output: %w", err)
	}

	return &Result{
		Rows: []ResultRow{
			{
				Title: fmt.Sprintf("Updated file %s", params.Filepath),
			},
		},
	}, nil
}

func (act *UpImageVersionAction) upVersion(img *image.Image, params Params) error {
	switch params.NewValue {
	case "major":
		img.Version = img.Version.UpMajor()
	case "minor":
		img.Version = img.Version.UpMinor()
	case "patch":
		img.Version = img.Version.UpPatch()
	default:
		return fmt.Errorf("unsupported version up method %q", params.NewValue)
	}
	return nil
}

func (act *UpImageVersionAction) dryRun(newContent string, params Params) (*Result, error) {
	return &Result{
		Rows: []ResultRow{
			{
				Title:   fmt.Sprintf("Dry Run for update %s", params.Pointer),
				Content: fmt.Sprintf("New content: \n%s", newContent),
			},
		},
	}, nil
}
