package evaluator

import (
	"github.com/ForeverSRC/paimon-interpreter/pkg/ast"
	"github.com/ForeverSRC/paimon-interpreter/pkg/object"
)

func Eval(node ast.Node) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalStatements(n.Statements)
	case *ast.ExpressionStatement:
		return Eval(n.Expression)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return getBoolean(n.Value)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func getBoolean(val bool) *object.Boolean {
	if val {
		return object.TRUE
	} else {
		return object.FALSE
	}
}
