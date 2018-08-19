package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
	"strings"
)

// ObjectType token.Tokenと同じような型
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
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

// Type fulfill the object.Object interface
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

// Inspect fulfill the object.Object interface
func (i *Integer) Inspect() string { return fmt.Sprintf("%d", i.Value) }

// Boolean boolean型
type Boolean struct {
	Value bool
}

// Type fulfill the object.Object interface
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

// Inspect fulfill the object.Object interface
func (b *Boolean) Inspect() string { return fmt.Sprintf("%t", b.Value) }

// String Go言語のString型をそのまま利用
type String struct {
	Value string
}

// Type fulfill the object.Object interface
func (s *String) Type() ObjectType { return STRING_OBJ }

// Inspect fulfill the object.Object interface
func (s *String) Inspect() string { return s.Value }

// Array 配列 Goの配列をそのまま利用
type Array struct {
	Elements []Object
}

// Type fulfill the object.Object interface
func (ao *Array) Type() ObjectType { return ARRAY_OBJ }

// Inspect fulfill the object.Object interface
func (ao *Array) Inspect() string {
	var out bytes.Buffer

	var elements []string
	for _, e := range ao.Elements {
		elements = append(elements, e.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}

// Null null型 nullを実装することでnullの実装について学ぶ
type Null struct {
}

// Type fulfill the object.Object interface
func (n *Null) Type() ObjectType { return NULL_OBJ }

// Inspect fulfill the object.Object interface
func (n *Null) Inspect() string { return "null" }

// ReturnValue Objectをラップする
type ReturnValue struct {
	Value Object
}

// Type fulfill the object.Object interface
func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }

// Inspect fulfill the object.Object interface
func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }

// Error Error Objectをラップする
type Error struct {
	Message string
}

// Type fulfill the object.Object interface
func (e *Error) Type() ObjectType { return ERROR_OBJ }

// Inspect fulfill the object.Object interface
func (e *Error) Inspect() string { return "ERROR: " + e.Message }

// Function function
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

// Type fulfill the object.Object interface
func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

// Inspect fulfill the object.Object interface
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

type BuiltinFunction func(args ...Object) Object

// Builtin ビルトイン関数
type Builtin struct {
	Fn BuiltinFunction
}

// Type fulfill the object.Object interface
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

// Inspect fulfill the object.Object interface
func (b *Builtin) Inspect() string { return "builtin function" }

// NewEnclosedEnvironment 外側の環境を参照する新しい環境を生成する
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

type HashKey struct {
	Type  ObjectType
	Value int64
}

// Hashable Hash化可能なオブジェクト
type Hashable interface {
	HashKey() HashKey
}

// HashKey boolean
func (b *Boolean) HashKey() HashKey {
	var value int64

	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: int64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: int64(h.Sum64())}
}

// HashPair HashKey生成元のKeyとValue
type HashPair struct {
	Key   Object
	Value Object
}

// Hash Monkey hash(map)　object
// Memo: map[HashKey]Objectにしない理由: Inspectでkey/valueを表示するため
// HashKeyには実際のkeyは格納されていない. (P242)
type Hash struct {
	Pairs map[HashKey]HashPair
}

// Type fulfill the object.Object interface
func (h *Hash) Type() ObjectType { return HASH_OBJ }

// Inspect fulfill the object.Object interface
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	var pairs []string
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s", pair.Key.Inspect(), pair.Value.Inspect()))
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
