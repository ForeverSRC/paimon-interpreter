package parser

import (
	"testing"

	"github.com/ForeverSRC/paimon-interpreter/pkg/lexer"
)

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a*b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a+b+c", "((a + b) + c)"},
		{"a+b-c", "((a + b) - c)"},
		{"a*b*c", "((a * b) * c)"},
		{"a+b/c", "(a + (b / c))"},
		{"a+b*c+d/e-f", "(((a + (b * c)) + (d / e)) - f)"},
		{"5<4==6<5", "((5 < 4) == (6 < 5))"},
		{"true", "true"},
		{"false", "false"},
		{"3>5==false", "((3 > 5) == false)"},
		{"3<5==false", "((3 < 5) == false)"},
		{"1+(1+2)+4", "((1 + (1 + 2)) + 4)"},
		{"(5+4)*2", "((5 + 4) * 2)"},
		{"-(5+5)", "(-(5 + 5))"},
		{"!(true==true)", "(!(true == true))"},
		{"a+add(b*c)+d", "((a + add((b * c))) + d)"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParseErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
