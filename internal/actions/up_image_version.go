package actions

import (
	"fmt"
	githuboutput "github.com/ci-space/github-output"
	"github.com/goccy/go-yaml/parser"
	"strings"

	gyaml "github.com/goccy/go-yaml"

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
	file, err := act.fs.ReadFile(params.Filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %v", err)
	}

	path, err := gyaml.PathString(preparePointer(params.Pointer))
	if err != nil {
		return nil, fmt.Errorf("failed to parse pointer: %v", err)
	}

	var imageString string

	err = path.Read(file.Reader(), &imageString)
	if err != nil {
		return nil, fmt.Errorf("failed to read yaml: %v", err)
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

	yFile, err := parser.ParseBytes(file.Content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	err = path.ReplaceWithReader(yFile, strings.NewReader(vImage.String()))
	if err != nil {
		return nil, err
	}

	newContent := yFile.String()

	if params.DryRun {
		return act.dryRun(newContent, params)
	}

	err = act.fs.WriteFile(params.Filepath, []byte(newContent))
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
