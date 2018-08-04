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
		t.Errorf("object is not integer. got=%T (%+v). expected=%d", evaluated, evaluated, expected)
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

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestEvalIfExpression(t *testing.T) {
	tests := [] struct {
		input    string
		expected bool
	}{
		{"if (true) { true } else {false}", true},
		{"if (true) { false }", false},
		{"if (1 < 2) { true } else {false}", true},
		{"if (1 == 2) { true } else {false}", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestReturnStatement(t *testing.T) {
	tests := [] struct {
		input    string
		expected int64
	}{
		{`
if (true) {
	return 10;
} else {
	return 20;
}
`, 10},
		{`
if (10 > 1) {
	if (10 > 1) {
		return 1;
    }
	return 10;
}
`, 1},
		{`if (1 > 0){10; return 9; 8}`, 9},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testBooleanObject(t *testing.T, evaluated object.Object, expected bool) bool {
	result, ok := evaluated.(*object.Boolean)
	if !ok {
		t.Errorf("evaluated is not a boolean. got=%T (%+v). expected=%t", evaluated, evaluated, expected)
		return false
	}

	if result.Value != expected {
		t.Errorf("evaluated result %t is not expected %v.", evaluated, expected)
		return false;
	}
	return true;
}

func TestNull(t *testing.T) {
	evaluated := testEval("null")
	if evaluated != nil {
		t.Errorf("evaluated is not null")
	}
}

func TestEvalInteger(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"1", 1},
		{"-10", -10},
		{"1 + 1", 2},
		{"10 - 3", 7},
		{"2 * 2 * 2", 8},
		{"9 / 3 * 2", 6},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(s string) object.Object {
	l := lexer.New(s)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}
