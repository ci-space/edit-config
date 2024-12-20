package markup

import "io"

type Document interface {
	Read(pointer string, value interface{}) error
	UpdateValue(pointer string, value io.Reader) error
	String() string
}
