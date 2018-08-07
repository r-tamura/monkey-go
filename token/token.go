package token

// TokenType string
// Tokenタイプ識別子
type TokenType string

// Token token
// トークンはトークンタイプとトークン値を持つ
type Token struct {
	Type    TokenType
	Literal string
}

// monkeyにおけるトークン一覧
const (
	ILLEGAL = "ILLEGAL" // 無効なトークン
	EOF     = "EOF"     // End of file

	// Identifier + literals
	IDENT = "IDENT"
	INT   = "INT"

	// Operators
	ASSIGN = "="
	PLUS   = "+"

	// Delimiters
	COMMA     = ","
	SEMICOLON = ";"

	LPAREN = "("
	RPAREN = ")"
	LBRACE = "{"
	RBRACE = "}"

	// Keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
)

// Builtin Identifier
var keywords = map[string]TokenType{
	"fn":  FUNCTION,
	"let": LET,
}

// LookupIdent ビルトインIdentifierであればビルトインのTokenTypeを返す
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
