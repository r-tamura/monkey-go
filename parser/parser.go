package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Parser Parser
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

// New Return new Parser instance
func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	// Read tow tokens, so curToken and peekToken are both set.
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseProgram 後で実装
// 繰り返しトークンを進めながら、現在のトークンタイプからどのようにパースするかを決定する
func (p *Parser) ParseProgram() *ast.Program {
	return nil
}
