package actions

import (
	"fmt"
	"github.com/artarts36/yamlpath"
)

type GetAction struct {
}

func NewGetAction() *GetAction {
	return &GetAction{}
}

func (g *GetAction) Run(content *yamlpath.Document, params Params) (*Result, error) {
	elem, err := content.Get(yamlpath.NewPointer(params.Pointer))
	if err != nil {
		return nil, fmt.Errorf("failed to get value from path: %w", err)
	}

	res, err := elem.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal value: %w", err)
	}

	return &Result{
		Rows: []ResultRow{
			{
				Title:   fmt.Sprintf("Value of %s", params.Pointer),
				Content: string(res),
			},
		},
	}, nil
}
