package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

// ObjectType token.Tokenと同じような型
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR_OBJ"
	FUNCTION_OBJ     = "FUNCTION_OBJ"
)

// Object 全ての値は異なる型で定義される
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Integer integer型
type Integer struct {
	Value int64
}

// Type fullfil the object.Object interface
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Inspect fullfil the object.Object interface
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean boolean型
type Boolean struct {
	Value bool
}

// Type fullfil the object.Object interface
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Inspect fullfil the object.Object interface
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// Null null型 nullを実装することでnullの実装について学ぶ
type Null struct {
}

// Type fullfil the object.Object interface
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Inspect fullfil the object.Object interface
func (n *Null) Inspect() string { return "null" }

// ReturnValue Objectをラップする
type ReturnValue struct {
	Value Object
}

// Type fullfil the object.Object interface
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Inspect fullfil the object.Object interface
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error Error Objectをラップする
type Error struct {
	Message string
}

// Type fullfil the object.Object interface
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Inspect fullfil the object.Object interface
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type fullfil the object.Object interface
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// Inspect fullfil the object.Object interface
func (f *Function) Inspect() string {
	// Memo: string結合よりもbytes.Bufferを使った方がパフォーマンスが良い
	// https://machiel.me/post/bytes-buffer-for-string-concatenation-in-go/
	var out bytes.Buffer

	var params []string
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ","))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// NewEnclosedenvironment 外側の環境を参照する新しい環境を生成する
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// NewEnvironment create new environment object
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// Environment 環境: 束縛されている変数の一覧を持つ
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Get get a value from the environment
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set set a value to the environment
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
