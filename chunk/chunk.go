package chunk

type Chunk struct {
	Code []uint8
}

func (c *Chunk) Write(code []uint8) (n int, err error) {
	c.Code = append(c.Code, code...)

	return len(code), nil
}

type OpCode uint8

const (
	OpReturn OpCode = iota
)
