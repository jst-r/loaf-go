package vm

import (
	"errors"
	"fmt"
	"io"
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
	err      error
	objects  *value.ObjPool
	globals  map[string]Value

	Stdout io.StringWriter
	Stderr io.StringWriter
}

func New() *VM {
	return &VM{Chunk: nil,
		ip: 0, stackTop: 0,
		objects: value.NewObjPool(),
		globals: make(map[string]Value),
		Stdout:  os.Stdout,
		Stderr:  os.Stderr,
	}
}

type InterpretResult int

const (
	InterpretOk InterpretResult = iota
	InterpretCompileError
	InterpretRuntimeError
)

func (v *VM) Interpret(chunk *bytecode.Chunk) InterpretResult {
	v.Chunk = chunk
	v.objects.Append(chunk.Objects)

	v.run()

	if v.err != nil {
		v.Stderr.WriteString(v.err.Error())
		return InterpretRuntimeError
	}

	return InterpretOk
}

func (v *VM) run() {
	defer func() {
		if r := recover(); r != nil {
			v.runtimeError(r.(string))
		}
	}()

	for {
		v.traceInstruction()

		switch v.readByte() {
		case bytecode.OpReturn:
			return
		case bytecode.OpConstant:
			constant := v.readConstant()
			v.push(constant)
		case bytecode.OpAdd:
			v.add()
		case bytecode.OpSubtract:
			v.binaryOp(subtract)
		case bytecode.OpMultiply:
			v.binaryOp(multiply)
		case bytecode.OpDivide:
			v.binaryOp(divide)
		case bytecode.OpNegate:
			v.negate()
		case bytecode.OpNil:
			v.push(value.Nil)
		case bytecode.OpTrue:
			v.push(value.Bool(true))
		case bytecode.OpFalse:
			v.push(value.Bool(false))
		case bytecode.OpNot:
			v.push(value.Bool(v.pop().IsFalsey()))
		case bytecode.OpEqual:
			b := v.pop()
			a := v.pop()
			v.push(value.Bool(value.ValuesEqual(a, b)))
		case bytecode.OpGreater:
			v.binaryOp(greater)
		case bytecode.OpLess:
			v.binaryOp(less)
		case bytecode.OpPop:
			v.pop()
		case bytecode.OpDefineGlobal:
			name := v.readString()
			v.globals[name] = v.peek(0)
			v.pop()
		case bytecode.OpGetGlobal:
			name := v.readString()
			value, ok := v.globals[name]
			if !ok {
				panic("Undefined variable: " + name)
			}
			v.push(value)
		case bytecode.OpSetGlobal:
			name := v.readString()
			_, ok := v.globals[name]
			if !ok {
				panic("Undefined variable: " + name)
			}
			v.globals[name] = v.peek(0)
		case bytecode.OpGetLocal:
			index := int(v.readByte())
			v.push(v.stack[index])
		case bytecode.OpSetLocal:
			index := int(v.readByte())
			v.stack[index] = v.peek(0)
		case bytecode.OpPrint:
			v.Stdout.WriteString(v.pop().FormatString())
			v.Stdout.WriteString("\n")
		}
	}
}

func (v *VM) runtimeError(msg string) {
	instruction := v.ip - 1
	line := v.Chunk.Lines.Find(instruction)
	errString := msg + "\n" + fmt.Sprintf("line %d in script\n", line)
	v.err = errors.New(errString)
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

func (v *VM) readString() string {
	return v.readConstant().AsString().Str
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
