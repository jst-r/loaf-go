package vm

import (
	"fmt"
	"os"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
)

type Value = value.Value

const StackSize = 1024

type VM struct {
	Chunk    *bytecode.Chunk
	ip       int
	stack    [StackSize]Value
	stackTop int
}

func New() *VM {
	return &VM{Chunk: nil, ip: 0, stackTop: 0}
}

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) Interpret(chunk *bytecode.Chunk) InterpretResult {
	v.Chunk = chunk

	return v.run()
}

func (v *VM) run() InterpretResult {
	for {
		v.traceInstruction()
		switch v.readByte() {
		case bytecode.OpReturn:
			fmt.Println(v.pop().FormatString())
			return InterpretOk
		case bytecode.OpConstant:
			constant := v.readConstant()
			v.push(constant)
		case bytecode.OpAdd:
			v.binaryOp(add)
		case bytecode.OpSubtract:
			v.binaryOp(subtract)
		case bytecode.OpMultiply:
			v.binaryOp(multiply)
		case bytecode.OpDivide:
			v.binaryOp(divide)
		case bytecode.OpNegate:
			if !v.peek(0).IsFloat() {
				v.runtimeError("Operand must be a number")
				return InterpretRuntimeError
			}
			v.push(value.Float(-v.pop().AsFloat()))
		}
	}
}

func (v *VM) runtimeError(format string, args ...interface{}) {
	os.Stderr.WriteString(fmt.Sprintf(format, args...) + "\n")
	instruction := v.ip - 1
	line := v.Chunk.Lines.Find(instruction)
	os.Stderr.WriteString(fmt.Sprintf("line %d in script\n", line))
}

func (v *VM) readByte() uint8 {
	b := v.Chunk.Code[v.ip]
	v.ip += 1
	return b
}

func (v *VM) readConstant() Value {
	index := int(v.readByte())
	return v.Chunk.Constants[index]
}

func (v *VM) push(value Value) {
	if v.stackTop >= StackSize {
		panic("stack overflow")
	}
	v.stack[v.stackTop] = value
	v.stackTop += 1
}

func (v *VM) pop() Value {
	if v.stackTop <= 0 {
		panic("stack underflow")
	}
	v.stackTop -= 1
	return v.stack[v.stackTop]
}

func (v *VM) peek(distance int) Value {
	if v.stackTop-distance-1 < 0 {
		panic("stack underflow")
	}
	return v.stack[v.stackTop-distance-1]
}

func (v *VM) binaryOp(op func(a, b Value) Value) {
	b := v.pop()
	a := v.pop()
	v.push(op(a, b))
}
