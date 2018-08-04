package evaluator

import (
	"../ast"
	"../object"
	"fmt"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, environment *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatement(node.Statements, environment)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, environment)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements, environment)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return convertNativeBooleanToObject(node.Value)
	case *ast.StringLiteral:
		return &object.String{Value: node.Value}
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, environment)
	case *ast.InfixExpression:
		return evalInfixExpression(node, environment)
	case *ast.IfExpression:
		return evalIfExpression(node, environment)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(node.ReturnValue, environment)}
	case *ast.LetStatement:
		val := Eval(node.Value, environment)
		if isError(val) {
			return val
		}
		environment.Set(node.Name.Value, val)
	case *ast.FunctionLiteral:
		return evalFunction(node, environment)
	case *ast.Identifier:
		return evalIdentifierExpression(node, environment)
	case *ast.CallExpression:
		return evalCallExpression(node, environment)
	}
	return nil
}
func evalCallExpression(callExpression *ast.CallExpression, environment *object.Environment) object.Object {
	function := Eval(callExpression.Function, environment)

	if (isError(function)) {
		return function
	}

	var parameters []object.Object

	for _, argmentExpression := range callExpression.Arguments {
		evaluated := Eval(argmentExpression, environment)
		if (isError(evaluated)) {
			return evaluated;
		}
		parameters = append(parameters, evaluated)
	}

	return applyFunction(function.(*object.Function), parameters, environment)
}
func applyFunction(function *object.Function, arguments []object.Object, environment *object.Environment) object.Object {
	enclosingEnvironment := object.NewEnclosingEnvironment(environment)

	for i, argument := range arguments {
		enclosingEnvironment.Set(function.Parameters[i].Value, argument)
	}

	return Eval(function.Body, enclosingEnvironment)
}
func evalFunction(literal *ast.FunctionLiteral, environment *object.Environment) object.Object {
	return &object.Function{
		Body:        literal.Body,
		Parameters:  literal.Parameters,
		Environment: environment,
	}
}
func evalIdentifierExpression(identifier *ast.Identifier, environment *object.Environment) object.Object {
	if identifier.Value == "null" {
		return NULL
	}
	value, ok := environment.Get(identifier.Value)
	if !ok {
		return newError("Identifier not found: %s", identifier)
	}
	return value
}
func evalBlockStatement(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, environment)

		switch result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		}
	}

	return result
}
func evalIfExpression(ifExpression *ast.IfExpression, environment *object.Environment) object.Object {
	condition := Eval(ifExpression.Condition, environment)

	var ret object.Object

	if condition == TRUE {
		ret = Eval(ifExpression.Consequence, environment)
	} else if ifExpression.Alternative != nil {
		ret = Eval(ifExpression.Alternative, environment)
	} else {
		ret = NULL
	}

	if isError(ret) {
		return ret
	}

	return ret
}
func evalInfixExpression(infixExpression *ast.InfixExpression, environment *object.Environment) object.Object {
	left := Eval(infixExpression.Left, environment)
	right := Eval(infixExpression.Right, environment)

	if isError(left) {
		return left
	}
	if isError(right) {
		return right
	}

	switch left.(type) {
	case *object.Integer:
		return evalInfixIntegerOperator(infixExpression.Operator, left, right)
	case *object.Boolean:
		return evalInfixBooleanOperator(infixExpression.Operator, left, right)
	case *object.String:
		return evalInfixStringOperator(infixExpression.Operator, left, right)
	default:
		return newError("Unsupported operator: %s %s %s", left.Type(), infixExpression.Operator, right.Type())
	}
}
func evalInfixStringOperator(operator string, left object.Object, right object.Object) object.Object {
	leftString, leftOk := left.(*object.String)
	rightString, rightOk := right.(*object.String)
	if !(leftOk && rightOk) {
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	switch operator {
	case "+":
		return &object.String{Value: leftString.Value + rightString.Value}
	default:
		return newError("Unsupported operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func isError(target object.Object) bool {
	_, ok := target.(*object.Error)
	return ok
}
func evalInfixBooleanOperator(operator string, left object.Object, right object.Object) object.Object {
	leftBoolean, leftOk := left.(*object.Boolean)
	rightBoolean, rightOk := right.(*object.Boolean)
	if !(leftOk && rightOk) {
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}

	leftValue := leftBoolean.Value
	rightValue := rightBoolean.Value

	switch operator {
	case "==":
		return convertNativeBooleanToObject(leftValue == rightValue)
	case "!=":
		return convertNativeBooleanToObject(leftValue != rightValue)
	default:
		return newError("Unsupported operator: %s %s %s", left.Type(), operator, right.Type())
	}
}
func evalInfixIntegerOperator(operator string, left object.Object, right object.Object) object.Object {
	leftInteger, leftOk := left.(*object.Integer)
	rightInteger, rightOk := right.(*object.Integer)
	if !(leftOk && rightOk) {
		return newError("Type mismatch: %s %s %s", left.Type(), operator, right.Type())
	}
	leftValue := leftInteger.Value
	rightValue := rightInteger.Value

	var result int64

	switch operator {
	case "+":
		result = leftValue + rightValue
	case "-":
		result = leftValue - rightValue
	case "*":
		result = leftValue * rightValue
	case "/":
		result = leftValue / rightValue
	case "<":
		return convertNativeBooleanToObject(leftValue < rightValue)
	case "<=":
		return convertNativeBooleanToObject(leftValue <= rightValue)
	case ">":
		return convertNativeBooleanToObject(leftValue > rightValue)
	case ">=":
		return convertNativeBooleanToObject(leftValue >= rightValue)
	case "==":
		return convertNativeBooleanToObject(leftValue == rightValue)
	case "!=":
		return convertNativeBooleanToObject(leftValue != rightValue)
	default:
		return newError("Unsupported operator: %s %s %s", left.Type(), operator, right.Type())
	}
	return &object.Integer{Value: result}
}

func evalPrefixExpression(prefixExpression *ast.PrefixExpression, environment *object.Environment) object.Object {
	right := Eval(prefixExpression.Right, environment)
	if isError(right) {
		return right
	}

	switch prefixExpression.Operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return newError("Unsupported operator: %s %s", prefixExpression.Operator, right.Type())
	}
}
func evalMinusOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	default:
		return newError("Unsupported operator: %s %s", "-", right.Type())
	}
}
func evalBangOperatorExpression(target object.Object) object.Object {
	switch target {
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

func convertNativeBooleanToObject(value bool) *object.Boolean {
	if value {
		return TRUE
	} else {
		return FALSE
	}
}
func evalStatement(statements []ast.Statement, environment *object.Environment) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement, environment)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}
func newError(message string, argumentTypes ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, argumentTypes...)}
}
