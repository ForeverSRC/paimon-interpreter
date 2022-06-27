package ast

import "github.com/ForeverSRC/paimon-interpreter/pkg/token"

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
