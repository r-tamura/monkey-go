package compiler

import (
	"monkey/ast"
	"monkey/code"
	"monkey/object"
)

// Compiler a Compler for Monkey Lnaguage
type Compiler struct {
	instructions code.Instructions
	constants    []object.Object
}

// New a constructor of the Compiler
func New() *Compiler {
	return &Compiler{
		instructions: code.Instructions{},
		constants:    []object.Object{},
	}
}

// Compile compile parsed AST
func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

// Bytecode ???
func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.instructions,
		Constants:    c.constants,
	}
}

// Bytecode ???
type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object // ???
}
