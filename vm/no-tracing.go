//go:build !trace_execution

package vm

import "github.com/jst-r/loaf-go/chunk"

func traceInstruction(ip int, chunk *chunk.Chunk) {
	// noop
}
