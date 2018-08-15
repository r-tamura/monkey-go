package object

import "fmt"

// ObjectType token.Tokenと同じような型
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
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
