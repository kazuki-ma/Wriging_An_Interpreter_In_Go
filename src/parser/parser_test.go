package parser

import "testing"
import "../ast"
import (
	"../lexer"
	"../token"
	"fmt"
	"log"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5", "x", 5},
		{"let y = true", "y", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {
		program := parseProgramWithParserErrors(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d",
				1, len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, letStmt.Name)
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `
		return 5;
		return 10;
		return 993322;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("Program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStatement, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if returnStatement.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not `return`, got %q",
				returnStatement.TokenLiteral())
			continue
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := `foobar;`

	program := parseProgramWithParserErrors(t, input)

	tests := []ast.Statement{
		&ast.ExpressionStatement{Token: token.Token{Type: token.IDENT, Literal: "foobar"},
			Expression: &ast.Identifier{Token: token.Token{Type: token.IDENT, Literal: "foobar"}, Value: "foobar"}},
	}

	testStatement(t, program, tests)
}

func TestPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
	}

	for _, tt := range prefixTests {
		program := parseProgramWithParserErrors(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statement. got=%d",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statemtents[0] is not *ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testStatement(t *testing.T, program *ast.Program, tests []ast.Statement) {
	if len(program.Statements) != len(tests) {
		t.Fatalf("Program has not enough statement. Got=%d",
			len(program.Statements))
	}

	for i := range program.Statements {
		if test, ok := tests[i].(*ast.ExpressionStatement); ok {
			stmt, ok := program.Statements[i].(*ast.ExpressionStatement)
			if !ok {
				t.Fatalf("program.Statement[%d] is not ast.ExpressionStatement. got=%T",
					i, program.Statements[i])
			}

			if expectedIdent, isIdent := test.Expression.(*ast.Identifier); isIdent {
				expression, ok := stmt.Expression.(*ast.Identifier)
				if !ok {
					t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
				}
				if expression.Value != "foobar" {
					t.Errorf("expression.Value not %s. got %s", expectedIdent.Value, expression.Value)
				}
				if expression.TokenLiteral() != "foobar" {
					t.Errorf("expression.TokenLiteral() not %s. got %s",
						expectedIdent.TokenLiteral(), expression.TokenLiteral())
				}
				continue
			}
			if expectedIntegerLiteral, isIntegerLiteral := test.Expression.(*ast.IntegerLiteral); isIntegerLiteral {
				expression, ok := stmt.Expression.(*ast.IntegerLiteral)
				if !ok {
					t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
				}
				if expression.Value != expectedIntegerLiteral.Value {
					t.Errorf("expression.Value not %d. got %d", expectedIntegerLiteral.Value, expression.Value)
				}
				if expression.TokenLiteral() != expectedIntegerLiteral.TokenLiteral() {
					t.Errorf("expression.TokenLiteral() not %s. got %s",
						expectedIntegerLiteral.TokenLiteral(), expression.TokenLiteral())
				}
				continue
			}
			t.Errorf("Not coverd in expression test. %T", test.Expression)
		}

	}
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.OperatorExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, int64(v))
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}
func testBooleanLiteral(t *testing.T, expression ast.Expression, value bool) bool {
	bo, ok := expression.(*ast.Boolean)

	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", expression)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%T", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
	program := parseProgramWithParserErrors(t, input)

	tests := []ast.Statement{
		&ast.ExpressionStatement{Token: token.Token{Type: token.INT, Literal: "5"},
			Expression: &ast.IntegerLiteral{Token: token.Token{Type: token.INT, Literal: "5"}, Value: 5}},
	}
	testStatement(t, program, tests)
}

func parseProgramWithParserErrors(t *testing.T, input string) *ast.Program {
	p := New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)
	return program
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		program := parseProgramWithParserErrors(t, tt.input)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not a pointer to ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("program.Statements[0].Expression is not compatible with *ast.InfixExpression. got=%T",
				stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

type IOPair struct {
	Input  string
	Output string
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []IOPair{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"fn(x){ x; }",
			"fn(x) {x;}",
		},
		{
			"fn(x, y){ x + y; }",
			"fn(x,y) {(x + y);}",
		},
	}

	testParsingUsingString(tests, t)
}

func TestGrouping(t *testing.T) {
	tests := []IOPair{
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			Input:  "- (5 + 5)",
			Output: "(-(5 + 5))",
		},
		{
			Input:  "!(true == true)",
			Output: "(!(true == true))",
		},
	}

	testParsingUsingString(tests, t)
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	program := parseProgramWithParserErrors(t, input)

	if len(program.Statements) != 1 {
		t.Errorf("program body does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression")
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got %+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	program := parseProgramWithParserErrors(t, input)

	if len(program.Statements) != 1 {
		t.Errorf("program body does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression")
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}
	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not *ast.ExpressionStatement. got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func testParsingUsingString(tests []IOPair, t *testing.T) {
	for _, tt := range tests {
		l := lexer.New(tt.Input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		log.Printf("%v ==> %v", tt.Input, program)

		actual := program.String()
		if actual != tt.Output {
			t.Errorf("expected=%q, got=%q", tt.Output, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"true", "true"},
		{"false", "false"},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
	}

	for _, tt := range tests {
		program := parseProgramWithParserErrors(t, tt.input)
		es := program.Statements[0].(*ast.ExpressionStatement)

		if tt.expected == es.String() {
			continue
		}

		t.Errorf("%s is not parsed as expected'%s'", es, tt.expected)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		program := parseProgramWithParserErrors(t, tt.input)
		exp := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.InfixExpression)

		if exp.Operator != tt.operator {
			t.Errorf("Operator expected %s. got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			continue
		}
		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			continue
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!true", "!", true},
		{"!false", "!", false},
	}

	for _, tt := range prefixTests {
		program := parseProgramWithParserErrors(t, tt.input)
		exp := program.Statements[0].(*ast.ExpressionStatement).Expression.(*ast.PrefixExpression)

		log.Printf("%s", exp)

		if exp.Operator != tt.operator {
			t.Errorf("Operator expected %s. got=%s", tt.operator, exp.Operator)
		}
		if !testLiteralExpression(t, exp.Right, tt.value) {
			continue
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y){ x + y; }`
	log.Printf("input: %s", input)

	program := parseProgramWithParserErrors(t, input)

	if len(program.Statements) != 1 {
		t.Errorf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Expression[0] is not *ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	functionLiteral, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not *ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	testLiteralExpression(t, functionLiteral.Parameters[0], "x")
	testLiteralExpression(t, functionLiteral.Parameters[1], "y")

	if len(functionLiteral.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d",
			len(functionLiteral.Body.Statements))
	}

	bodyStatement, ok := functionLiteral.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			functionLiteral.Body.Statements[0])
	}

	testInfixExpression(t, bodyStatement.Expression, "x", "+", "y")
}

func TestStringLiteralParsing(t *testing.T) {
	inputs := `"TEST"`

	program := parseProgramWithParserErrors(t, inputs)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements has wrong length. Expected=%d. Got=%d",
			1, len(program.Statements))
	}

	expressionStatement := program.Statements[0].(*ast.ExpressionStatement)
	log.Printf("%s", expressionStatement)
}

func TestCallExpressionParsing(t *testing.T) {
	input := `add(1, 2 * 3, 4 + 5)`
	program := parseProgramWithParserErrors(t, input)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. Got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not an *ast.CallExpression. got=%T",
			stmt.Expression)
	}
	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	return true
}
