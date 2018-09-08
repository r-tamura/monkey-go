package code

import (
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
)

type Definition struct {
	Name          string
	OperandWidths []int
}

var definitions = map[Opcode]*Definition{
	OpConstant: {"OpConstand", []int{2}},
}

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
		}
		offset += width
	}
	return instruction
}
