package bytecode

import (
	"fmt"
	"strings"
)

func (c *Chunk) Disassemble(name string) string {
	return c.NewDisassembler(name).disassemble()
}

type disassembler struct {
	*Chunk
	name    string
	offset  int
	builder *strings.Builder
}

func (c *Chunk) NewDisassembler(name string) *disassembler {
	return &disassembler{c, name, 0, &strings.Builder{}}
}

func (d *disassembler) SetOffset(offset int) {
	d.offset = offset
}

func (d *disassembler) String() string {
	return d.builder.String()
}

func (d *disassembler) disassemble() string {
	d.builder.WriteString("=== ")
	d.builder.WriteString(d.name)
	d.builder.WriteString(" ===\n")

	for d.offset < len(d.Code) {
		d.builder.WriteString(fmt.Sprintf("%04d ", d.offset))
		if d.offset > 0 && d.Lines.Find(d.offset-1) == d.Lines.Find(d.offset) {
			d.builder.WriteString("   | ")
		} else {
			d.builder.WriteString(fmt.Sprintf("%4d ", d.Lines.Find(d.offset)))
		}
		d.DisassembleInstruction()
	}
	return d.builder.String()
}

func (d *disassembler) DisassembleInstruction() {
	switch d.Code[d.offset] {
	case OpReturn:
		d.simpleInstruction("OP_RETURN")
	case OpConstant:
		d.constantInstruction("OP_CONSTANT")
	case OpSetGlobal:
		d.constantInstruction("OP_DEFINE_GLOBAL")
	case OpGetGlobal:
		d.constantInstruction("OP_GET_GLOBAL")
	case OpNot:
		d.simpleInstruction("OP_NOT")
	case OpNegate:
		d.simpleInstruction("OP_NEGATE")
	case OpAdd:
		d.simpleInstruction("OP_ADD")
	case OpSubtract:
		d.simpleInstruction("OP_SUBTRACT")
	case OpMultiply:
		d.simpleInstruction("OP_MULTIPLY")
	case OpDivide:
		d.simpleInstruction("OP_DIVIDE")
	case OpNil:
		d.simpleInstruction("OP_NIL")
	case OpTrue:
		d.simpleInstruction("OP_TRUE")
	case OpFalse:
		d.simpleInstruction("OP_FALSE")
	case OpEqual:
		d.simpleInstruction("OP_EQUAL")
	case OpGreater:
		d.simpleInstruction("OP_GREATER")
	case OpLess:
		d.simpleInstruction("OP_LESS")
	case OpPrint:
		d.simpleInstruction("OP_PRINT")
	case OpPop:
		d.simpleInstruction("OP_POP")
	default:
		d.builder.WriteString(fmt.Sprintf("unknown instruction %d\n", d.Code[d.offset]))
		d.offset += 1
	}
}

func (d *disassembler) simpleInstruction(name string) {
	d.builder.WriteString(name)
	d.builder.WriteString("\n")
	d.offset += 1
}

func (d *disassembler) constantInstruction(name string) {
	index := int(d.Code[d.offset+1])
	value := d.Constants[index].FormatString()
	d.builder.WriteString(fmt.Sprintf("%-16s %4d %s\n", name, index, value))
	d.offset += 2
}
