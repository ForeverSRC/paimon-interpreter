package ast

import "github.com/ForeverSRC/paimon-interpreter/pkg/token"

// Program AST的根结点
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

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
