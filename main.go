package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/value"
	"github.com/jst-r/loaf-go/vm"
)

func main() {
	vm := vm.New()

	prog := bytecode.Chunk{}
	ind := prog.AddConstant(value.Float(1.0))

	prog.Write(bytecode.OpConstant, 1)
	prog.Write(uint8(ind), 1)
	prog.Write(bytecode.OpNegate, 1)
	prog.Write(bytecode.OpReturn, 2)
	fmt.Println(prog.Disassemble("main"))

	res := vm.Interpret(&prog)
	fmt.Println("Interpret result:", res)
}
