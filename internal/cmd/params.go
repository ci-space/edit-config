package cmd

import "github.com/ci-space/edit-config/internal/actions"

type Params struct {
	Filepath string
	Action   actions.Name
	Pointer  string
	NewValue string
	DryRun   bool
}
