package ast

import (
	"bytes"

	"github.com/ForeverSRC/paimon-interpreter/pkg/token"
)

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	var out bytes.Buffer

	out.WriteString("{\n")
	for _, s := range b.Statements {
		out.WriteString("  " + s.String())
	}
	out.WriteString("\n}")

	return out.String()
}

func (b *BlockStatement) statementNode() {
}
