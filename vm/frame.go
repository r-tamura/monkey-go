package vm

import (
	"monkey/code"
	"monkey/object"
)

// Frame short for call frame and stack frame (activation record)
// CompiledFunctionとip(instruction point)をプロパティに持つ
// monkeyの実装では関数のみにFrameを使用する
// basePointer: stack上でのフレームの開始位置を記憶する
type Frame struct {
	cl          *object.Closure
	ip          int
	basePointer int
}

// NewFrame 新しいFrameを生成する
func NewFrame(cl *object.Closure, basePointer int) *Frame {
	return &Frame{cl: cl, ip: -1, basePointer: basePointer}
}

// Instructions Frameが参照しているCompileFunctionのInstructionsを取得する
func (f *Frame) Instructions() code.Instructions {
	return f.cl.Fn.Instructions
}
