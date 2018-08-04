package evaluator

import (
	"../ast"
	"../object"
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

		if returnResult, ok := result.(*object.ReturnValue); ok {
			return returnResult
		}
	}

	return result
}
func evalIfExpression(ifExpression *ast.IfExpression) object.Object {
	condition := Eval(ifExpression.Condition)

	if condition == TRUE {
		return Eval(ifExpression.Consequence)
	}

	if ifExpression.Alternative != nil {
		return Eval(ifExpression.Alternative)
	}

	return NULL
}
func evalInfixExpression(infixExpression *ast.InfixExpression) object.Object {
	left := Eval(infixExpression.Left)
	right := Eval(infixExpression.Right)

	switch left.(type) {
	case *object.Integer:
		return evalInfixIntegerOperator(infixExpression.Operator, left, right)
	case *object.Boolean:
		return evalInfixBooleanOperator(infixExpression.Operator, left, right)
	default:
		return NULL
	}
}
func evalInfixBooleanOperator(operator string, left object.Object, right object.Object) object.Object {
	if left.(*object.Boolean) == nil {
		return NULL
	}
	if right.(*object.Boolean) == nil {
		return NULL
	}
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	switch operator {
	case "==":
		return convertNativeBooleanToObject(leftValue == rightValue)
	case "!=":
		return convertNativeBooleanToObject(leftValue != rightValue)
	default:
		return NULL
	}
}
func evalInfixIntegerOperator(operator string, left object.Object, right object.Object) object.Object {
	if left.(*object.Integer) == nil {
		return NULL
	}
	if right.(*object.Integer) == nil {
		return NULL
	}
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

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
		return NULL
	}
	return &object.Integer{Value: result}
}

func evalPrefixExpression(prefixExpression *ast.PrefixExpression) object.Object {
	right := Eval(prefixExpression.Right)
	switch prefixExpression.Operator {
	case "!":
		return evalBandOperatorExpression(right)
	case "-":
		return evalMinusOperatorExpression(right)
	default:
		return NULL
	}
}
func evalMinusOperatorExpression(right object.Object) object.Object {
	switch right := right.(type) {
	case *object.Integer:
		return &object.Integer{Value: -right.Value}
	default:
		return NULL
	}

}
func evalBandOperatorExpression(target object.Object) object.Object {
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

		if returnResult, ok := result.(*object.ReturnValue); ok {
			return returnResult.Value
		}
	}

	return result
}
