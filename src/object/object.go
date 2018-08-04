package object

import (
	"ast"
	"bytes"
	"fmt"
	"strings"
)

type ObjectType string

const (
	INTEGER_OBJ      ObjectType = "INTEGER"
	BOOLEAN_OBJ      ObjectType = "BOOLEAN"
	NULL_OBJ         ObjectType = "NULL"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	ERROR_OBJ        ObjectType = "ERROR"
	FUNCTION_OBJ     ObjectType = "FUNCTION"
	STRING_OBJ       ObjectType = "STRING"
	BUILDIN_OBJ      ObjectType = "BUILDIN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}
func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}
func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType {
	return STRING_OBJ
}
func (s *String) Inspect() string {
	return fmt.Sprintf(`"%s"`, s.Value)
}

type Null struct {
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}
func (n *Null) Inspect() string {
	return "NULL"
}

type ReturnValue struct {
	Value Object
}

func (rValue *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ

}
func (rValue *ReturnValue) Inspect() string {
	return fmt.Sprintf("return %s", rValue.Inspect())
}

type Error struct {
	Message string
}

func (error *Error) Type() ObjectType {
	return ERROR_OBJ
}
func (error *Error) Inspect() string {
	return "ERROR: " + error.Message
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}
func NewEnclosingEnvironment(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]Object),
		outer: outer,
	}
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func (env *Environment) Get(name string) (Object, bool) {
	object, ok := env.store[name]
	if !ok && env.outer != nil {
		object, ok = env.outer.Get(name)
	}
	return object, ok
}
func (env *Environment) Set(name string, object Object) Object {
	env.store[name] = object
	return object
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment *Environment
}

func (function *Function) Type() ObjectType {
	return FUNCTION_OBJ
}
func (function *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, param := range function.Parameters {
		params = append(params, param.String())
	}

	out.WriteString("fn (")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(function.Body.String())
	out.WriteString("\n}")

	return out.String()
}

type BuildinFunction func(args ...Object) Object
type Buildin struct {
	Fn BuildinFunction
}

func (buildin *Buildin) Type() ObjectType {
	return BUILDIN_OBJ
}
func (buildin *Buildin) Inspect() string {
	return "buildin function"
}
