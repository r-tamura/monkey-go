package lexer

import (
	"testing"

	"monkey/token"
)

func TestNextToken(t *testing.T) {
	input := `=+{},;`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
	}

	l := New(input)
	for i, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expectedType {
			// %qは空白などをエスケープシーケンスに置き換えシングルクオートで囲んだ文字を表示する
			t.Fatalf("text[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("text[%d] - literal wrong. expected=%q, got%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}
