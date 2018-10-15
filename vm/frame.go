package vm

import (
	"monkey/code"
	"monkey/object"
)

// Frame short for call frame and stack frame (activation record)
// CompiledFunctionとip(instruction point)をプロパティに持つ
// monkeyの実装では関数のみにFrameを使用する
type Frame struct {
	fn *object.CompiledFunction
	ip int
}

// NewFrame 新しいFrameを生成する
func NewFrame(fn *object.CompiledFunction) *Frame {
	return &Frame{fn: fn, ip: -1}
}

// Instructions Frameが参照しているCompileFunctionのInstructionsを取得する
func (f *Frame) Instructions() code.Instructions {
	return f.fn.Instructions
}
