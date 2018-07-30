package evaluator

import "testing"
import "../lexer"
import "../parser"
import "../object"

func TestEvalIntegerExpression(t *testing.T) {
	tests := [] struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}
func testIntegerObject(t *testing.T, evaluated interface{}, expected int64) bool {
	result, ok := evaluated.(*object.Integer)
	if !ok {
		t.Errorf("object is not integer. got=%T (%+v)", evaluated, evaluated)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}
	return true
}

func testEval(s string) object.Object {
	l := lexer.New(s)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}
