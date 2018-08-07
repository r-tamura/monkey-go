package lexer

import "monkey/token"

// Lexer 文字列をトークンへ変換する
type Lexer struct {
	input        string
	position     int  // 現在見ている文字の位置
	readPosition int  // 現在読んでいる文字の位置(positionの後)
	ch           byte // 現在見ている文字 ASCIIのみを扱う仕様(TODO: Unicode対応)
}

// New Lexerインスタンスを生成する
func New(input string) *Lexer {
	// Golangの仕様: 指定しないプロパティはZero valueが設定される
	// int => 0
	// string => ""
	// byte => 0
	l := &Lexer{input: input}

	// 現在位置を1文字目に設定
	l.readChar()
	return l
}

// NextToken 次のトークンを返します
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.SkipWhiteSpaces()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = token.EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isNumber(l.ch) {
			tok.Literal = l.readNumber()
			tok.Type = token.INT
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		// 入力文字列数を超える場合は終了(NULL文字とする)
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isNumber(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

func isLetter(ch byte) bool {
	// _を変数名に使用できる
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// SkipWhiteSpaces 空白スキップ
func (l *Lexer) SkipWhiteSpaces() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}
