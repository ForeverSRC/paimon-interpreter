package evaluator

import (
	"fmt"

	"github.com/ForeverSRC/paimon-interpreter/pkg/ast"
	"github.com/ForeverSRC/paimon-interpreter/pkg/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch n := node.(type) {
	case *ast.Program:
		return evalProgram(n, env)
	case *ast.ExpressionStatement:
		return Eval(n.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: n.Value}
	case *ast.Boolean:
		return getBoolean(n.Value)
	case *ast.PrefixExpression:
		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(n.Operator, right)
	case *ast.InfixExpression:
		left := Eval(n.Left, env)
		if isError(left) {
			return left
		}

		right := Eval(n.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(n.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(n, env)
	case *ast.IfExpression:
		return evalIfExpression(n, env)
	case *ast.ReturnStatement:
		val := Eval(n.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}

	case *ast.LetStatement:
		val := Eval(n.Value, env)
		if isError(val) {
			return val
		}

		env.Set(n.Name.Value, val)
	case *ast.Identifier:
		return evalIdentifier(n, env)
	case *ast.FunctionLiteral:
		params := n.Parameters
		body := n.Body
		return &object.Function{Parameters: params, Body: body, Env: env}
	case *ast.CallExpression:
		function := Eval(n.Function, env)
		if isError(function) {
			return function
		}

		args := evalExpressions(n.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}

		return applyFunction(function, args)
	}

	return nil
}

func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range program.Statements {
		result = Eval(stmt, env)

		switch res := result.(type) {
		case *object.ReturnValue:
			return res.Value
		case *object.Error:
			return res
		}

	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValObj || rt == object.ErrorObj {
				return result
			}
		}
	}

	return result
}

func getBoolean(val bool) *object.Boolean {
	if val {
		return TRUE
	} else {
		return FALSE
	}
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: %s%s", operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerObj {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() != right.Type():
		return newError("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	case left.Type() == object.IntegerObj && right.Type() == object.IntegerObj:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return getBoolean(left == right)
	case operator == "!=":
		return getBoolean(left != right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	lv := left.(*object.Integer).Value
	rv := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: lv + rv}
	case "-":
		return &object.Integer{Value: lv - rv}
	case "*":
		return &object.Integer{Value: lv * rv}
	case "/":
		return &object.Integer{Value: lv / rv}
	case "<":
		return getBoolean(lv < rv)
	case ">":
		return getBoolean(lv > rv)
	case "==":
		return getBoolean(lv == rv)
	case "!=":
		return getBoolean(lv != rv)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ie *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(ie.Condition, env)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ie.Consequence, env)
	} else if ie.Alternative != nil {
		return Eval(ie.Alternative, env)
	} else {
		return NULL
	}
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	val, ok := env.Get(node.Value)
	if !ok {
		return newError("identifier not found: %s", node.Value)
	}

	return val
}

func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	for _, e := range exps {
		evaluated := Eval(e, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}

		result = append(result, evaluated)
	}

	return result
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	function, ok := fn.(*object.Function)
	if !ok {
		return newError("not a function: %s", fn.Type())
	}

	extendEnv := extendFunctionEnv(function, args)
	evaluated := Eval(function.Body, extendEnv)
	return unwrapReturnValue(evaluated)
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)
	for paramIdx, param := range fn.Parameters {
		env.Set(param.Value, args[paramIdx])
	}

	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if rtVal, ok := obj.(*object.ReturnValue); ok {
		return rtVal.Value
	}

	return obj
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{
		Message: fmt.Sprintf(format, a...),
	}
}

func isError(obj object.Object) bool {
	return obj != nil && obj.Type() == object.ErrorObj
}
