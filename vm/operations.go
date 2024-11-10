package vm

import "github.com/jst-r/loaf-go/value"

func (v *VM) binaryOp(op func(a, b Value) Value) {
	b := v.pop()
	a := v.pop()
	v.push(op(a, b))
}

func add(a, b Value) Value {
	return value.Float(a.AsFloat() + b.AsFloat())
}

func subtract(a, b Value) Value {
	return value.Float(a.AsFloat() - b.AsFloat())
}

func multiply(a, b Value) Value {
	return value.Float(a.AsFloat() * b.AsFloat())
}

func divide(a, b Value) Value {
	return value.Float(a.AsFloat() / b.AsFloat())
}
