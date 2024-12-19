package cmd

import (
	"fmt"
	cli "github.com/artarts36/singlecli"
	"github.com/artarts36/yamlpath"
	"github.com/ci-space/edit-config/internal/actions"
	"github.com/ci-space/edit-config/internal/fs"
	"os"
)

type Command struct {
	output cli.Output
}

func NewCommand(output cli.Output) *Command {
	return &Command{output: output}
}

func (c *Command) Run(params ActionParams) error {
	acts := actions.CreateActions(fs.NewLocal())
	act, ok := acts[params.Action]
	if !ok {
		return fmt.Errorf("unknown action: %s", params.Action)
	}

	contentBytes, err := os.ReadFile(params.Filepath)
	if err != nil {
		return fmt.Errorf("failed to read file %q: %w", params.Filepath, err)
	}

	content, err := yamlpath.Unmarshal(contentBytes)
	if err != nil {
		return fmt.Errorf("failed to parse file %q: %w", params.Filepath, err)
	}

	return c.runAction(act, content, params)
}

func (c *Command) runAction(act actions.Action, doc *yamlpath.Document, params ActionParams) error {
	res, err := act.Run(doc, actions.Params{
		Filepath: params.Filepath,
		Pointer:  params.Pointer,
		NewValue: params.NewValue,
		DryRun:   params.DryRun,
	})
	if err != nil {
		return err
	}

	for _, row := range res.Rows {
		fmt.Println(row.Title)

		if row.Content != "" {
			fmt.Println()
			fmt.Println(row.Content)
		}
	}

	return nil
}
