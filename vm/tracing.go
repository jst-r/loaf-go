//go:build trace_execution

package vm

import (
	"fmt"
)

func (v *VM) traceInstruction() {
	fmt.Print("        ")
	for i := v.stackTop - 1; i >= 0; i-- {
		fmt.Printf("[ %s ]", v.stack[i].String())
	}
	fmt.Println()
	dis := v.Chunk.NewDisassembler("main")
	dis.SetOffset(v.ip)
	dis.DisassembleInstruction()
	fmt.Print(dis.String())
}
