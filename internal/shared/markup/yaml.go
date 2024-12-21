package markup

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/goccy/go-yaml"
	"github.com/goccy/go-yaml/ast"
	"github.com/goccy/go-yaml/parser"
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

func (y *YAMLDocument) Append(pointer string, value interface{}) error {
	ptr, err := y.getPointer(pointer)
	if err != nil {
		return fmt.Errorf("could not get pointer for %q: %v", pointer, err)
	}

	node, err := y.readPathNode(ptr, y.file)
	if err != nil {
		return fmt.Errorf("could not read node for %q: %v", pointer, err)
	}

	switch n := node.(type) {
	case *ast.StringNode:
		valueStr, ok := value.(string)
		if !ok {
			return fmt.Errorf("value for %q must be a string", pointer)
		}

		n.Value += valueStr
		n.Token.Value += valueStr

		return ptr.ReplaceWithNode(y.file, n)
	case *ast.IntegerNode:
		valueInt, ok := value.(int)
		if !ok {
			return fmt.Errorf("value for %q must be an int", pointer)
		}

		switch currVal := n.Value.(type) {
		case int64:
			n.Value = currVal + int64(valueInt)
			n.Token.Value = fmt.Sprintf("%d", n.Value)
		case uint64:
			n.Value = currVal + uint64(valueInt)
			n.Token.Value = fmt.Sprintf("%d", n.Value)
		default:
			return fmt.Errorf("value for %q must be an integer", pointer)
		}

		n.Value = valueInt

		return ptr.ReplaceWithNode(y.file, n)
	case *ast.SequenceNode:
		value, err = y.prepareNewValueForSequence(pointer, n, value)
		if err != nil {
			return fmt.Errorf("could not prepare new value for %q: %v", pointer, err)
		}

		marshalled, merr := yaml.Marshal(value)
		if merr != nil {
			return merr
		}

		astFile, perr := parser.ParseBytes(marshalled, parser.ParseComments)
		if perr != nil {
			return fmt.Errorf("could not parse %q: %v", pointer, perr)
		}

		n.Values = append(n.Values, astFile.Docs[0])

		return ptr.ReplaceWithNode(y.file, n)
	default:
		return fmt.Errorf("value for %q must be a string, int or sequence", pointer)
	}
}

func (y *YAMLDocument) UpdateValue(pointer string, value interface{}) error {
	ptr, err := y.getPointer(pointer)
	if err != nil {
		return fmt.Errorf("could not get pointer for %q: %v", pointer, err)
	}

	switch v := value.(type) {
	case string:
		return ptr.ReplaceWithReader(y.file, strings.NewReader(v))
	case []byte:
		return ptr.ReplaceWithReader(y.file, bytes.NewReader(v))
	default:
		marshalled, merr := yaml.Marshal(v)
		if merr != nil {
			return merr
		}

		return ptr.ReplaceWithReader(y.file, bytes.NewReader(marshalled))
	}
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

func (y *YAMLDocument) readPathNode(p *yaml.Path, r io.Reader) (ast.Node, error) {
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		return nil, err
	}
	f, err := parser.ParseBytes(buf.Bytes(), parser.ParseComments)
	if err != nil {
		return nil, err
	}
	node, err := p.FilterFile(f)
	if err != nil {
		return nil, err
	}
	return node, nil
}

func (y *YAMLDocument) prepareNewValueForSequence(
	pointer string,
	seq *ast.SequenceNode,
	value interface{},
) (interface{}, error) {
	if len(seq.Values) == 0 {
		return value, nil
	}

	var newValue interface{}
	var err error

	switch seq.Values[0].(type) {
	case *ast.IntegerNode:
		switch v := value.(type) {
		case int:
			newValue = v
		case string:
			newValue, err = strconv.Atoi(v)
			if err != nil {
				return nil, fmt.Errorf("value for %q must be an int: %w", pointer, err)
			}
		}
	case *ast.BoolNode:
		switch v := value.(type) {
		case bool:
			newValue = v
		case string:
			newValue, err = strconv.ParseBool(v)
			if err != nil {
				return nil, fmt.Errorf("value for %q must be an int: %w", pointer, err)
			}
		}
	case *ast.FloatNode:
		switch v := value.(type) {
		case float64:
			newValue = v
		case string:
			newValue, err = strconv.ParseFloat(v, 64)
			if err != nil {
				return nil, fmt.Errorf("value for %q must be an float: %w", pointer, err)
			}
		}
	case *ast.StringNode:
		switch v := value.(type) {
		case string:
			newValue = v
		default:
			return nil, fmt.Errorf("value for %q must be a string, got %T", pointer, value)
		}
	}

	return newValue, nil
}
