package object

import (
	"bytes"
	"strings"

	"github.com/ForeverSRC/paimon-interpreter/pkg/ast"
)

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() Type {
	return FuncObj
}

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := make([]string, 0, len(f.Parameters))

	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString("){\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
