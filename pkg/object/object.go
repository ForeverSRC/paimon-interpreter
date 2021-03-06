package object

type Type string

type Object interface {
	Type() Type
	Inspect() string
}

const (
	IntegerObj   Type = "Integer"
	BooleanOjb   Type = "Boolean"
	NullObj      Type = "Null"
	ReturnValObj Type = "ReturnValue"
	ErrorObj     Type = "Error"
	FuncObj      Type = "Function"
)
