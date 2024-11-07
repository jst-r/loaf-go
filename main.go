package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
)

func main() {
	prog := chunk.Chunk{}
	ind := prog.AddConstant(1.0)

	prog.Write(chunk.OpConstant, 1)
	prog.Write(uint8(ind), 1)
	prog.Write(chunk.OpReturn, 1)
	prog.Write(chunk.OpReturn, 2)
	fmt.Println(prog.Disassemble("main"))
}
