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
		if d.offset > 0 && d.lines.Find(d.offset-1) == d.lines.Find(d.offset) {
			d.builder.WriteString("   | ")
		} else {
			d.builder.WriteString(fmt.Sprintf("%4d ", d.lines.Find(d.offset)))
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
	case OpNegate:
		d.simpleInstruction("OP_NEGATE")
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
	value := d.Constants[index].String()
	d.builder.WriteString(fmt.Sprintf("%-16s %4d %s\n", name, index, value))
	d.offset += 2
}
