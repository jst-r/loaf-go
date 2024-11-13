package bytecode

import (
	"fmt"
	"log"

	"github.com/jst-r/loaf-go/value"
)

type Chunk struct {
	Code      []uint8
	Lines     LineInfo
	Constants value.ValueArray
	Objects   *value.ObjPool
}

func NewChunk() *Chunk {
	return &Chunk{Objects: value.NewObjPool()}
}

func (c *Chunk) Write(code uint8, line int) error {
	c.Code = append(c.Code, code)
	c.Lines.AddLine(len(c.Code)-1, line)

	return nil
}

func (c *Chunk) WriteSlice(code []uint8, lines []int) error {
	if len(lines) != len(code) {
		return fmt.Errorf("invalid code length")
	}

	for i, line := range lines { // doesn't look too got but it's linear
		c.Lines.AddLine(len(c.Code)+i-1, line)
	}
	c.Code = append(c.Code, code...)

	return nil
}

func (c *Chunk) AddConstant(value value.Value) (index int) {
	c.Constants = append(c.Constants, value)
	index = len(c.Constants) - 1
	if index > 255 {
		log.Fatal("too many constants")
	}
	return index
}

// Assumes that the line numbers are monotonically increasing
type LineInfo struct {
	spans []LineSpan
}

// Start and end are inclusive
type LineSpan struct {
	Line        int
	StartOffset int
	EndOffset   int
}

func (l *LineInfo) AddLine(offset int, line int) {
	if len(l.spans) == 0 {
		l.spans = append(l.spans, LineSpan{line, offset, offset})
	} else if line == l.spans[len(l.spans)-1].Line {
		l.spans[len(l.spans)-1].EndOffset = offset
	} else {
		l.spans = append(l.spans, LineSpan{line, offset, offset})
	}
}

// This results in quadratic time, but it should be fine
func (l *LineInfo) Find(offset int) int {
	for _, span := range l.spans {
		if offset >= span.StartOffset && offset <= span.EndOffset {
			return span.Line
		}
	}
	return -1
}
