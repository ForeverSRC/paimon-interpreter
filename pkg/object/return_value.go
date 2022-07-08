package object

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() Type {
	return ReturnValObj
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}
