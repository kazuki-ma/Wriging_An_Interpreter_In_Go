package evaluator

import "testing"
import "../lexer"
import "../parser"
import (
	"../object"
	"log"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
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
	tests := []struct {
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
	tests := []struct {
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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1 + true", "ERROR: Type mismatch: INTEGER + BOOLEAN"},
		{"5 + true; 5;", "ERROR: Type mismatch: INTEGER + BOOLEAN"},
		{"-true", "ERROR: Unsupported operator: - BOOLEAN"},
		{"true + false", "ERROR: Unsupported operator: BOOLEAN + BOOLEAN"},
		{"5; true + false; 5;", "ERROR: Unsupported operator: BOOLEAN + BOOLEAN"},
		{"if (true){ true + false; }", "ERROR: Unsupported operator: BOOLEAN + BOOLEAN"},
		{"true + true + true", "ERROR: Unsupported operator: BOOLEAN + BOOLEAN"},
		{"foobar", "ERROR: Identifier not found: foobar"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errorObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("Error object expected but got=%s (%T)\nCase: %s", evaluated, evaluated, tt.input)
			continue
		}

		if errorObj.Inspect() != tt.expected {
			t.Errorf("Error message is not equal to expected value.\n  Got=%s\n  Expected=%s\n  Case=%s",
				errorObj.Inspect(), tt.expected, tt.input)
			continue
		}
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
		return false
	}
	return true
}

func TestNull(t *testing.T) {
	evaluated := testEval("null")

	if evaluated == nil {
		t.Errorf("evaluated is not null. Got=%s (%T)", evaluated, evaluated)
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
	return Eval(program, object.NewEnvironment())
}

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 5; a;", 5},
		{"let a = 5 * 5; a;", 25},
		{"let a = 5; let b = a; b", 5},
		{"let a = 5; let b = 10; let c = a + b; c", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := `fn (x) { x + 2; }`

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Errorf("evaluated is not function. got=%s(%T)", evaluated, evaluated)
		return
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("Parameter length mismatch. Expected=%d, got=%d", 1, len(fn.Parameters))
	}
	parameter := fn.Parameters[0]
	if parameter.Value != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", parameter)
	}

	expectedBody := "{(x + 2);}"

	if fn.Body.String() != expectedBody {
		expectError(t, expectedBody, fn.Body.String())
	}
}

func expectError(t *testing.T, expected interface{}, actual interface{}) {
	t.Errorf("Expected %s (%T) but got %s (%T)", expected, expected, actual, actual)
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"fn(x){x}(1);", 1},
		{"let a = fn(x){x}; a(1);", 1},
		{"let b = 1; let a = fn(){b}; a();", 1},
	}

	for _, tt := range tests {
		log.Printf("----- Case %s", tt.input)
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestClosure(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let a = 10; let f = fn(){a}; f()", 10},
	}

	for _, tt := range tests {
		log.Printf("----- Case: %s", tt.input)

		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestString(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"TEST"`, "TEST"},
		{`fn(){"TEST"}()`, "TEST"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		if evaluated == nil {
			t.Errorf("Got nil for case = %s", tt.input)
			continue
		}
		if evaluated.Type() != "STRING" {
			t.Errorf("Type is not %s. got=%s", object.STRING_OBJ, evaluated.Type())
			continue
		}
	}
}
