package ast

import "github.com/ForeverSRC/paimon-interpreter/pkg/token"

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
