package actions

import "github.com/artarts36/yamlpath"

type Action interface {
	Run(content *yamlpath.Document, params Params) (*Result, error)
}

type Params struct {
	Filepath string
	Pointer  string
	NewValue string
	DryRun   bool
}

type Result struct {
	Rows []ResultRow
}

type ResultRow struct {
	Title   string
	Content string
}
