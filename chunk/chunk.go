package chunk

type Chunk struct {
	Code      []uint8
	Constants ValueArray
}

func (c *Chunk) Write(code []uint8) (n int, err error) {
	c.Code = append(c.Code, code...)

	return len(code), nil
}

func (c *Chunk) AddConstant(value Value) (index int) {
	c.Constants = append(c.Constants, value)
	return len(c.Constants) - 1
}

const (
	OpConstant uint8 = iota
	OpReturn
)
