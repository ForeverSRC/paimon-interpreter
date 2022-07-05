package object

import "fmt"

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() Type {
	return BooleanOjb
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)