package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func TestLetstatement(t *testing.T) {
	input := `
		let x = 5;
		let y = 10;
		ley foobar = 838383;
	`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements, got %d.", len(program.Statements))
	}
	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
		return false
	}

	_, ok := s.(*ast.LetStatement)
	if !ok {
		// %Tは型を返す
		// %T	a Go-syntax representation of the type of the value
		t.Errorf("s not *ast.Letstatement. got=%T", s)
	}

	return true
}
