package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/compiler"
	"github.com/jst-r/loaf-go/vm"
)

var source = //
`var x = 1;
print x;
`

func main() {
	vm := vm.New()

	prog, errs := compiler.Compile(source)
	if len(errs) > 0 {
		fmt.Println("Compile errors:", errs)
		return
	}
	fmt.Println(prog.Disassemble("main"))
	res := vm.Interpret(prog)
	fmt.Println("Interpret result:", res)
}
