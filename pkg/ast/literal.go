package ast

import "github.com/ForeverSRC/paimon-interpreter/pkg/token"

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}

func (il *IntegerLiteral) String() string {
	return il.TokenLiteral()
}

func (il *IntegerLiteral) expressionNode() {
}
