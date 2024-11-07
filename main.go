package main

import (
	"fmt"

	"github.com/jst-r/loaf-go/chunk"
)

func main() {
	prog := chunk.Chunk{}
	prog.Write([]uint8{1, 2, 3})
	fmt.Println(prog.Code)
}
