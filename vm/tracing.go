//go:build trace_execution

package vm

import (
	"fmt"
)

func (v *VM) traceInstruction() {
	fmt.Print("        ")
	for i := v.ip - 1; i >= 0; i-- {
		fmt.Printf("[ %s ]", v.stack[i].String())
	}
	dis := v.Chunk.NewDisassembler("main")
	dis.SetOffset(v.ip)
	dis.DisassembleInstruction()
	fmt.Print(dis.String())
}
