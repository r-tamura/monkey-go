package compiler

type SymbolScope string

const (
	LocalScope  SymbolScope = "LOCAL"
	GlobalScope SymbolScope = "GLOBAL"
)

// Symbol プログラム内で宣言された変数を表す
type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

// SymbolTable Symbol表
type SymbolTable struct {
	Outer *SymbolTable
	// key: identifier / value: symbol
	store          map[string]Symbol
	numDefinitions int
}

// NewSymbolTable GlobalSymbolTable用
func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	return &SymbolTable{store: s}
}

// NewEnclosedSymbolTable LocalSymbolTable用
func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

// Define SymbolTableへSymbolを定義する
func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

// Resolve SymbolTableからSymbolを解決する
// 未定義Symbole名の場合は第二返り値がfalseとなる
func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		return s.Outer.Resolve(name)
	}
	return obj, ok
}
