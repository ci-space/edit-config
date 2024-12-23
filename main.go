package main

import (
	"context"
	cli "github.com/artarts36/singlecli"
	"github.com/ci-space/edit-config/internal/actions"
	"github.com/ci-space/edit-config/internal/cmd"
)

var (
	Version   = "v0.1.0"
	BuildDate = "2024-12-20 03:09:15"
)

func main() {
	app := &cli.App{
		BuildInfo: &cli.BuildInfo{
			Name:        "edit-config",
			Description: "edit-config is a tool that edits YAML files",
			Version:     Version,
			BuildDate:   BuildDate,
		},
		Action: run,
		Args: []*cli.ArgDefinition{
			{
				Name:        "file",
				Description: "path to YAML file",
				Required:    true,
			},
			{
				Name:        "action",
				Description: "action to edit",
				Required:    true,
				ValuesEnum:  actions.Names,
			},
			{
				Name:        "pointer",
				Description: "pointer to element in YAML file",
				Required:    true,
			},
			{
				Name:        "new-value",
				Description: "new value in YAML file",
			},
		},
		Opts: []*cli.OptDefinition{
			{
				Name: "dry-run",
			},
			{
				Name:        "separator",
				Description: "separate value",
				WithValue:   true,
			},
		},
	}

	app.RunWithGlobalArgs(context.Background())
}

func run(ctx *cli.Context) error {
	command := cmd.NewCommand(ctx.Output)

	action, err := actions.NameFromString(ctx.GetArg("action"))
	if err != nil {
		return err
	}

	return command.Run(cmd.Params{
		Filepath:  ctx.GetArg("file"),
		Action:    action,
		Pointer:   ctx.GetArg("pointer"),
		NewValue:  ctx.GetArg("new-value"),
		DryRun:    ctx.HasOpt("dry-run"),
		Separator: ctx.Opts["separator"],
	})
}
