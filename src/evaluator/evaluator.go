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

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatement(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.BlockStatement:
		return evalBlockStatement(node.Statements)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return convertNativeBooleanToObject(node.Value)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node)
	case *ast.InfixExpression:
		return evalInfixExpression(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		return &object.ReturnValue{Value: Eval(node.ReturnValue)}
	}
	return nil
}
func evalBlockStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		switch result.(type) {
		case *object.ReturnValue:
			return result
		case *object.Error:
			return result
		}
	}

	return result
}
func evalIfExpression(ifExpression *ast.IfExpression) object.Object {
	condition := Eval(ifExpression.Condition)

	var ret object.Object

	if condition == TRUE {
		ret = Eval(ifExpression.Consequence)
	} else if ifExpression.Alternative != nil {
		ret = Eval(ifExpression.Alternative)
	} else {
		ret = NULL
	}

	if isError(ret) {
		return ret
	}

	return ret
}
func evalInfixExpression(infixExpression *ast.InfixExpression) object.Object {
	left := Eval(infixExpression.Left)
	right := Eval(infixExpression.Right)

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
	default:
		return newError("Unsupported operator: %s %s %s", left.Type(), infixExpression.Operator, right.Type())
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

func evalPrefixExpression(prefixExpression *ast.PrefixExpression) object.Object {
	right := Eval(prefixExpression.Right)
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
		return FALSE;
	}
}
func evalStatement(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}
func newError(message string, argumentTypes ... interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, argumentTypes...)}
}
