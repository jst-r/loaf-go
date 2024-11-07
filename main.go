package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
	"github.com/jst-r/loaf-go/value"
	"github.com/jst-r/loaf-go/vm"
)

func main() {
	vm := vm.New()

	prog := chunk.Chunk{}
	ind := prog.AddConstant(value.Float(1.0))

	prog.Write(chunk.OpConstant, 1)
	prog.Write(uint8(ind), 1)
	prog.Write(chunk.OpReturn, 1)
	fmt.Println(prog.Disassemble("main"))

	res := vm.Interpret(&prog)
	fmt.Println("Interpret result:", res)
}
