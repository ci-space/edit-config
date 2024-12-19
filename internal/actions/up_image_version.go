package actions

import (
	"fmt"
	githuboutput "github.com/ci-space/github-output"

	"github.com/artarts36/yamlpath"

	"github.com/ci-space/edit-config/internal/fs"
	"github.com/ci-space/edit-config/internal/shared/image"
)

type UpImageVersionAction struct {
	fs fs.Filesystem
}

func NewUpImageVersionAction(fs fs.Filesystem) *UpImageVersionAction {
	return &UpImageVersionAction{fs: fs}
}

func (act *UpImageVersionAction) Run(doc *yamlpath.Document, params Params) (*Result, error) {
	img, err := doc.Get(yamlpath.NewPointer(params.Pointer))
	if err != nil {
		return nil, err
	}

	imageVal, err := img.AsScalar()
	if err != nil {
		return nil, err
	}

	imageString, ok := imageVal.(string)
	if !ok {
		return nil, fmt.Errorf("expected img to be a string")
	}

	vImage, err := image.ParseImage(imageString)
	if err != nil {
		return nil, err
	}

	switch params.NewValue {
	case "major":
		vImage.Version.UpMajor()
	case "minor":
		vImage.Version.UpMinor()
	case "patch":
		vImage.Version.UpPatch()
	default:
		return nil, fmt.Errorf("unsupported version up method %q", params.NewValue)
	}

	err = img.Update(nil, vImage.String())
	if err != nil {
		return nil, fmt.Errorf("failed to update image version: %w", err)
	}

	newContent, err := doc.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal new document content: %w", err)
	}

	if params.DryRun {
		return act.dryRun(newContent, params)
	}

	err = act.fs.WriteFile(params.Filepath, newContent)
	if err != nil {
		return nil, fmt.Errorf("failed to write updated config: %w", err)
	}

	err = githuboutput.WhenAvailable(func() error {
		return githuboutput.Write("new-value", vImage.String())
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

func (act *UpImageVersionAction) dryRun(newContent []byte, params Params) (*Result, error) {
	return &Result{
		Rows: []ResultRow{
			{
				Title:   fmt.Sprintf("Dry Run for update %s", params.Pointer),
				Content: fmt.Sprintf("New content: \n%s", newContent),
			},
		},
	}, nil
}
