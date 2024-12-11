package main

import (
	"fmt"
	"os"

	"github.com/jst-r/loaf-go/compiler"
	"github.com/jst-r/loaf-go/vm"
)

var source = //
`var x = 1;
if (x == 1) {
	print "x is 1";
}
if (x == 2) {
	print "x is 2";
}
`

func main() {
	vm := vm.New()

	prog, errs := compiler.Compile(source)
	if len(errs) > 0 {
		fmt.Println("Compile errors:")
		for _, err := range errs {
			fmt.Println(err)
		}
		os.Exit(1)
	}
	fmt.Println(prog.Disassemble("main"))
	res := vm.Interpret(prog)
	fmt.Println("Interpret result:", res)
}
