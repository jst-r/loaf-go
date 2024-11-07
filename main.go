package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
)

func main() {
	prog := chunk.Chunk{}
	prog.Write([]uint8{chunk.OpReturn})
	fmt.Println(prog.Disassemble("main"))
}
