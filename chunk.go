package loafgo

type Chunk struct {
	code []uint8
}

func (c *Chunk) Write(code []uint8) (n int, err error) {
	c.code = append(c.code, code...)

	return len(code), nil
}

type OpCode uint8

const (
	OpReturn OpCode = iota
)
