package ast

import (
	"bytes"

	"github.com/ForeverSRC/paimon-interpreter/pkg/token"
)

// ReturnStatement return <表达式>
type ReturnStatement struct {
	// Token token.RETURN
	Token       token.Token
	ReturnValue Expression
}

func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

func (r *ReturnStatement) statementNode() {
}

func (r *ReturnStatement) String() string {
	var out bytes.Buffer

	out.WriteString(r.TokenLiteral() + " ")
	if r.ReturnValue != nil {
		out.WriteString(r.ReturnValue.String())
	}

	out.WriteString(";")

	return out.String()
}
