package chunk

import "fmt"

type Chunk struct {
	Code      []uint8
	lines     []int
	Constants ValueArray
}

func (c *Chunk) Write(code uint8, line int) error {
	c.Code = append(c.Code, code)
	c.lines = append(c.lines, line) // TODO: use RLE

	return nil
}

func (c *Chunk) WriteSlice(code []uint8, lines []int) error {
	if len(lines) != len(code) {
		return fmt.Errorf("invalid code length")
	}

	c.Code = append(c.Code, code...)
	c.lines = append(c.lines, lines...) // TODO: use RLE

	return nil
}

func (c *Chunk) AddConstant(value Value) (index int) {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

const (
	OpConstant uint8 = iota
	OpReturn
)
