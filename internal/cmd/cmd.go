package cmd

import (
	"fmt"
	cli "github.com/artarts36/singlecli"
	"github.com/ci-space/edit-config/internal/actions"
	"github.com/ci-space/edit-config/internal/fs"
)

type Command struct {
	output cli.Output
}

func NewCommand(output cli.Output) *Command {
	return &Command{output: output}
}

func (c *Command) Run(params Params) error {
	acts := actions.CreateActions(fs.NewLocal())
	act, ok := acts[params.Action]
	if !ok {
		return fmt.Errorf("unknown action: %s", params.Action)
	}

	return c.runAction(act, params)
}

func (c *Command) runAction(act actions.Action, params Params) error {
	res, err := act.Run(actions.Params{
		Filepath:  params.Filepath,
		Pointer:   params.Pointer,
		NewValue:  params.NewValue,
		DryRun:    params.DryRun,
		Separator: params.Separator,
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
