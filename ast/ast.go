package ast

import "monkey/token"

// Node 全てのASTの要素はNodeである
type Node interface {
	// デバッグテスト用
	TokenLiteral() string
}

// Statement ステートメントNode
type Statement interface {
	Node
	// ダミーメソッドStatementとExpressionをコンパイラが区別できるように
	statementNode()
}

// Expression 式Node
type Expression interface {
	Node
	// ダミーメソッドStatementとExpressionをコンパイラが区別できるように
	expressionNode()
}

// Program 全てのASTのroot
type Program struct {
	Statements []Statement
}

// TokenLiteral Statementの先頭のTokenListerlを返す
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

// LetStatement let句
type LetStatement struct {
	Token token.Token // the token.LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}

// TokenLiteral Nodeリテラル実装
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }

type ReturnStatement struct {
	Token       token.Token // the token.RETURN token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode()       {}
func (rs *ReturnStatement) TokenLiteral() string { return rs.Token.Literal }

// Identifier 識別子
// let句ではIdentifierは値を返さないが、Monkey内で位置によっては値を生成する場合があるのでExpressionを実装(P45)
type Identifier struct {
	Token token.Token // the token.IDENT token
	Value string
}

func (i *Identifier) expressionNode() {}

// TokenLiteral Nodeリテラル実装
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
