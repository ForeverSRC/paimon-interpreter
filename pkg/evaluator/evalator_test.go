package evaluator

import (
	"testing"

	"github.com/ForeverSRC/paimon-interpreter/pkg/lexer"
	"github.com/ForeverSRC/paimon-interpreter/pkg/object"
	"github.com/ForeverSRC/paimon-interpreter/pkg/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-10", -10},
		{"-5", -5},
		{"5+5", 10},
		{"5-5", 0},
		{"5*5", 25},
		{"5/5", 1},
		{"5+5-5*5", -15},
		{"5-5+5", 5},
		{"2+(5+4)-3", 8},
		{"12/4+3*2-10", -1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1<2", true},
		{"1>2", false},
		{"1==1", true},
		{"1!=1", false},
		{"1==2", false},
		{"1!=2", true},
		{"(1<2)==true", true},
		{"(1>2)==false", true},
		{"((5>5)==true)==false", true},
	}

	for index, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, index, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, index int, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("index[%d] object is not Boolean. got=%T (%+v)", index, obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("index[%d] object has wrong value. got=%t, want=%t", index, result.Value, expected)
		return false
	}

	return true
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for index, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, index, evaluated, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if(true){10}", 10},
		{"if(false){10}", nil},
		{"if(1){10}", 10},
		{"if(1<2){10}", 10},
		{"if(1>2){10}", nil},
		{"if(1>2){10}else{20}", 20},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"return 10;", 10},
		{"return 11;9;", 11},
		{"return 3*5;9;", 15},
		{"9;return 4*5;9;", 20},
		{`if(10>1){
					if(10>1){return 10;}
					return 1;
				}`, 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"5+true;", "type mismatch: Integer + Boolean"},
		{"5+true;5;", "type mismatch: Integer + Boolean"},
		{"-true", "unknown operator: -Boolean"},
		{"true+false", "unknown operator: Boolean + Boolean"},
		{"if(10>1){true+false}", "unknown operator: Boolean + Boolean"},
		{`if(10>1){
					if(10>1){return true+false;}
					return 1;
				}`, "unknown operator: Boolean + Boolean"},
		{"foobar", "identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expected {
			t.Errorf("wrong error message. expected=%q, got=%q", tt.expected, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a=5;a;", 5},
		{"let a=5*5;a;", 25},
		{"let a=5;let b=a;b;", 5},
		{"let a=5;let b=a;let c=a+b+5;c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x){x+2;}"
	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T(%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong numbers of parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0].String())
	}

	expectedBody := "{\n  (x + 2)\n}"
	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let ident=fn(x){x;};ident(5);", 5},
		{"let ident=fn(x){return x;};ident(10);", 10},
		{"let double=fn(x){x*2;};double(4);", 8},
		{"let add=fn(x,y){x+y;};add(5,4);", 9},
		{"let add=fn(x,y){x+y;};add(5+5,add(1,3));", 14},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
