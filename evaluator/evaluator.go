package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

var (
	// NULL Null用オブジェクト
	NULL = &object.Null{}
	// TRUE boolean用オブジェクト全てのtrueはこのオブジェクトとして評価される
	TRUE = &object.Boolean{Value: true}
	// FALSE boolean用オブジェクト全てのfalseはこのオブジェクトとして評価される
	FALSE = &object.Boolean{Value: false}
)

// Eval parserによって生成されたASTを評価する
func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	// Statements
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)

	// Expressions
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		// booleanは2パターンしか存在しないので同じオブジェクトを参照するように (P146)
		return nativeBoolToBooleanObject(node.Value)
	}

	return nil
}

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range stmts {
		result = Eval(statement)
	}

	return result
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	}
	return FALSE
}
