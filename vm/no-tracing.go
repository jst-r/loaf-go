//go:build !trace_execution

package vm

func (v *VM) traceInstruction() {
	// noop
}
