package vm

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
)

type VM struct {
	Chunk *chunk.Chunk
	ip    int
}

func New() *VM {
	return &VM{nil, 0}
}

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) Interpret(chunk *chunk.Chunk) InterpretResult {
	v.Chunk = chunk

	return v.run()
}

func (v *VM) run() InterpretResult {
	for {
		switch v.readByte() {
		case chunk.OpReturn:
			return InterpretOk
		case chunk.OpConstant:
			val := v.readConstant()
			fmt.Println("constant:", val)
		}
	}

	return InterpretRuntimeError
}

func (v *VM) readByte() uint8 {
	b := v.Chunk.Code[v.ip]
	v.ip += 1
	return b
}

func (v *VM) readConstant() chunk.Value {
	index := int(v.readByte())
	return v.Chunk.Constants[index]
}
