package ast

import (
	"bytes"

	"github.com/ForeverSRC/paimon-interpreter/pkg/token"
)

// LetStatement let <标识符> = <表达式>;
type LetStatement struct {
	// Token token.LET 语法单元
	Token token.Token
	// Name 保存标识符
	Name *Identifier
	// Value 产生值的表达式
	Value Expression
}

func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) statementNode() {
}

func (ls *LetStatement) String() string {
	var out bytes.Buffer

	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}

	out.WriteString(";")

	return out.String()
}

type Identifier struct {
	// Token token.IDENT 语法单元
	Token token.Token
	Value string
}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) expressionNode() {
}

func (i *Identifier) String() string {
	return i.Value
}
