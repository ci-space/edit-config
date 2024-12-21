package markup

type Document interface {
	Read(pointer string, value interface{}) error
	Append(pointer string, value interface{}) error
	UpdateValue(pointer string, value interface{}) error
	String() string
}
