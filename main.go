package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/bytecode"
	"github.com/jst-r/loaf-go/compiler"
	"github.com/jst-r/loaf-go/value"
	"github.com/jst-r/loaf-go/vm"
)

func main() {
	vm := vm.New()

	prog := bytecode.Chunk{}
	ind1 := prog.AddConstant(value.Float(1.0))
	ind2 := prog.AddConstant(value.Float(2.0))

	bcode := []uint8{
		bytecode.OpConstant,
		uint8(ind1),
		bytecode.OpConstant,
		uint8(ind2),
		bytecode.OpAdd,
		bytecode.OpNegate,
		bytecode.OpReturn,
	}
	prog.WriteSlice(bcode, make([]int, len(bcode)))

	// fmt.Println(prog.Disassemble("main"))

	res := vm.Interpret(&prog)
	fmt.Println("Interpret result:", res)

	scanner := compiler.NewScanner("(-1 + 2) * 3 - -4")
	for {
		token := scanner.Scan()
		fmt.Printf("%+v\n", token)
		if token.Type == compiler.TokenEof {
			break
		}
	}

	p, errs := compiler.Compile("(-1 + 2) * 3 - -4")
	if len(errs) > 0 {
		fmt.Println(errs)
		return
	}
	fmt.Println(p.Disassemble("main"))
}
