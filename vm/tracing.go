//go:build trace_execution

package vm

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
)

func traceInstruction(ip int, chunk *chunk.Chunk) {
	dis := chunk.NewDisassembler("main")
	dis.SetOffset(ip)
	dis.DisassembleInstruction()
	fmt.Print(dis.String())
}
