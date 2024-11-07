package chunk

import (
	"fmt"
	"strings"
)

func (c *Chunk) Disassemble(name string) string {
	builder := strings.Builder{}

	builder.WriteString("=== ")
	builder.WriteString(name)
	builder.WriteString(" ===\n")

	for offset := 0; offset < len(c.Code); {
		builder.WriteString(fmt.Sprintf("%04d ", offset))
		offset = c.disassembleInstruction(offset, &builder)
	}
	return builder.String()
}

func (c *Chunk) disassembleInstruction(offset int, builder *strings.Builder) (newOffset int) {
	switch c.Code[offset] {
	case OpReturn:
		return simpleInstruction("OP_RETURN", offset, builder)
	case OpConstant:
		return constantInstruction("OP_CONSTANT", offset, builder, c)
	default:
		builder.WriteString(fmt.Sprintf("unknown instruction %d\n", c.Code[offset]))
		return offset + 1
	}

}

func simpleInstruction(name string, offset int, builder *strings.Builder) int {
	builder.WriteString(name)
	builder.WriteString("\n")
	return offset + 1
}

func constantInstruction(name string, offset int, builder *strings.Builder, c *Chunk) int {
	index := int(c.Code[offset+1])
	value := c.Constants[index]
	builder.WriteString(fmt.Sprintf("%-16s %4d %f\n", name, index, value))
	return offset + 2
}
