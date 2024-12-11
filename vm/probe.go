package vm

type VMProbe struct {
	*VM
}

func (v *VM) Probe() *VMProbe {
	return &VMProbe{v}
}

func (v *VMProbe) StackTop() int {
	return v.stackTop
}

func (v *VMProbe) Stack() []Value {
	return v.stack[:v.stackTop]
}
