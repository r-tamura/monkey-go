package code

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Instructions Instruction consisft of an opcode and an optional number of operands.
type Instructions []byte

// Opcode Operator Code
type Opcode byte

const (
	// OpConstant OpConstrant
	OpConstant Opcode = iota
	// OpAdd Stack上から2つのデータを取り出し、加算した結果をStack上へ追加する
	OpAdd
	OpSub
	OpMul
	OpDiv
	// OpPop Stack上から1つのデータを取り出す
	OpPop

	// Boolean
	OpTrue
	OpFalse

	// Comparison Operators
	// Memo: less thanオペレータはgreater thanの左辺・右辺をコンパイラ上で入れ替えるだけなので定義しない
	// The expression 3 < 5 can be reorderd to 5 > 3 without changing its result.
	OpEqual
	OpNotEqual
	OpGreaterThan

	// Prefix Operators
	OpMinus
	OpBang

	// Jump
	OpJumpNotTruthy
	OpJump

	// Bindings
	OpGetGlobal
	OpSetGlobal

	// Null
	OpNull

	// Array
	// VMはstackからN個の要素をpop、各要素を評価して、 Object.ArrayStackへpushする
	OpArray
	// Hash
	OpHash
	// Index Operator
	OpIndex

	// Function
	OpCall
	OpReturnValue // Stackの最後の要素を返す
	OpReturn      // Nullを返す(返す値が存在しないFunctionで利用する)

	// Local Bindings
	OpGetLocal
	OpSetLocal

	OpGetBuiltin

	// Closure
	OpClosure
	OpGetFree
)

// Definition a defition of monkey instructions
type Definition struct {
	Name          string
	OperandWidths []int // Operandのバイト数
}

// Stack上から引数を取り出すOperatorは引数がないことがある
var definitions = map[Opcode]*Definition{
	OpConstant:      {"OpConstant", []int{2}},
	OpAdd:           {"OpAdd", []int{}},
	OpSub:           {"OpSub", []int{}},
	OpMul:           {"OpMul", []int{}},
	OpDiv:           {"OpDiv", []int{}},
	OpPop:           {"OpPop", []int{}},
	OpTrue:          {"OpTrue", []int{}},
	OpFalse:         {"OpFalse", []int{}},
	OpEqual:         {"OpEqual", []int{}},
	OpNotEqual:      {"OpNotEqual", []int{}},
	OpGreaterThan:   {"OpGreaterThan", []int{}},
	OpMinus:         {"OpMinus", []int{}},
	OpBang:          {"OpBang", []int{}},
	OpJumpNotTruthy: {"OpJumpNotTruthy", []int{2}},
	OpJump:          {"OpJump", []int{2}},
	OpGetGlobal:     {"OpGetGlobal", []int{2}},
	OpSetGlobal:     {"OpSetGlobal", []int{2}},
	OpNull:          {"OpNull", []int{}},
	OpArray:         {"OpArray", []int{2}},
	OpHash:          {"OpHash", []int{2}},
	OpIndex:         {"OpIndex", []int{}},
	OpCall:          {"OpCall", []int{1}},
	OpReturnValue:   {"OpReturnValue", []int{}},
	OpReturn:        {"OpReturn", []int{}},
	OpGetLocal:      {"OpGetLocal", []int{1}},
	OpSetLocal:      {"OpSetLocal", []int{1}},
	OpGetBuiltin:    {"OpGetBuiltin", []int{1}},
	// Closure
	// 2byte Constant pool上の関数の位置
	// 1byte 自由変数の数
	OpClosure: {"OpClosure", []int{2, 1}},
	OpGetFree: {"OpGetFree", []int{1}},
}

// Lookup Lookup
func Lookup(op byte) (*Definition, error) {
	def, ok := definitions[Opcode(op)]
	if !ok {
		return nil, fmt.Errorf("opcode %d undefined", op)
	}
	return def, nil
}

// Make make a single bytecode instruction
func Make(op Opcode, operands ...int) []byte {
	def, ok := definitions[op]
	if !ok {
		return []byte{}
	}

	instructionLen := 1
	for _, w := range def.OperandWidths {
		instructionLen += w
	}

	// instructionの領域を確保
	// instruction: [Opcode Operand1 Operand2 ...]
	instruction := make([]byte, instructionLen)
	// Opcodeを配置(1 byte)
	instruction[0] = byte(op)

	// Operandをそれぞれ配置
	// Operandのデータサイズにより配置方法を決定する
	offset := 1
	for i, o := range operands {
		width := def.OperandWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instruction[offset:], uint16(o))
		case 1:
			instruction[offset] = byte(o)
		}
		offset += width
	}
	return instruction
}

// String バイト列をdisassembleする
func (ins Instructions) String() string {
	var out bytes.Buffer

	i := 0
	for i < len(ins) {
		def, err := Lookup(ins[i])
		if err != nil {
			fmt.Fprintf(&out, "ERROR: %s^n", err)
			continue
		}
		// Opcodeに対応する型定義とInstructionsのオペランド部分
		operands, read := ReadOperands(def, ins[i+1:])
		fmt.Fprintf(&out, "%04d %s\n", i, ins.fmtInstruction(def, operands))

		// Opcoe("1" byte) + Operand length("read" byte)
		i += 1 + read
	}

	return out.String()
}

func (ins Instructions) fmtInstruction(def *Definition, operands []int) string {
	operandCount := len(def.OperandWidths)
	if len(operands) != operandCount {
		return fmt.Sprintf("ERROR: operand len %d does not math defined %d\n", len(operands), operandCount)
	}

	switch operandCount {
	case 0:
		return def.Name
	case 1:
		return fmt.Sprintf("%s %d", def.Name, operands[0])
	case 2:
		return fmt.Sprintf("%s %d %d", def.Name, operands[0], operands[1])
	}

	return fmt.Sprintf("ERROR: unhandled operandCount for %s\n", def.Name)
}

// ReadOperands Makeの反対の機能 InstructionsのOperand部分 -> Operandを返す
func ReadOperands(def *Definition, ins Instructions) ([]int, int) {
	operands := make([]int, len(def.OperandWidths))
	offset := 0
	for i, width := range def.OperandWidths {
		switch width {
		case 2: // オペランドバイト数が 2 byte
			operands[i] = int(ReadUint16(ins[offset:]))
		case 1: // オペランドバイト数が 1 byte
			operands[i] = int(ReadUint8(ins[offset:]))
		}
		offset += width
	}
	return operands, offset
}

// ReadUint16 VMから利用するために関数として分離
func ReadUint16(ins Instructions) uint16 {
	return binary.BigEndian.Uint16(ins)
}

// ReadUint8 VMから利用するために関数として分離
func ReadUint8(ins Instructions) uint8 { return uint8(ins[0]) }
