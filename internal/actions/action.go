package actions

type Action interface {
	Run(params Params) (*Result, error)
}

type Params struct {
	Filepath  string
	Pointer   string
	NewValue  string
	DryRun    bool
	Separator string
}

type Result struct {
	Rows []ResultRow
}

type ResultRow struct {
	Title   string
	Content string
}
