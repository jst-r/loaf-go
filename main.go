package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
)

func main() {
	prog := chunk.Chunk{}
	ind := prog.AddConstant(1.0)
	prog.Write([]uint8{chunk.OpReturn, chunk.OpConstant, uint8(ind)})
	fmt.Println(prog.Disassemble("main"))
}
