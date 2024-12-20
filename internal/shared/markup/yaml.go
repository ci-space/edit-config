package markup

import (
	"fmt"
	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
	"io"
	"strings"
)

type YAMLDocument struct {
	file *ast.File

	pointers map[string]*yaml.Path
}

func LoadYAMLDocument(content []byte) (*YAMLDocument, error) {
	file, err := parser.ParseBytes(content, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	return &YAMLDocument{file: file, pointers: map[string]*yaml.Path{}}, nil
}

func (y *YAMLDocument) Read(pointer string, value interface{}) error {
	ptr, err := y.getPointer(pointer)
	if err != nil {
		return fmt.Errorf("could not get pointer for %q: %v", pointer, err)
	}

	return ptr.Read(y.file, value)
}

func (y *YAMLDocument) UpdateValue(pointer string, value io.Reader) error {
	ptr, err := y.getPointer(pointer)
	if err != nil {
		return fmt.Errorf("could not get pointer for %q: %v", pointer, err)
	}

	return ptr.ReplaceWithReader(y.file, value)
}

func (y *YAMLDocument) String() string {
	return y.file.String()
}

func (y *YAMLDocument) getPointer(pointer string) (*yaml.Path, error) {
	ptr, ok := y.pointers[pointer]
	if !ok {
		var err error
		ptr, err = yaml.PathString(y.preparePointer(pointer))
		if err != nil {
			return nil, err
		}

		y.pointers[pointer] = ptr
	}

	return ptr, nil
}

func (y *YAMLDocument) preparePointer(pointer string) string {
	if pointer == "" {
		return pointer
	}

	if !strings.HasPrefix(pointer, "$.") {
		return "$." + pointer
	}

	return pointer
}
