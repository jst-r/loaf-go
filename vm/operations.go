package vm

import "github.com/jst-r/loaf-go/value"

func (v *VM) binaryOp(op func(a, b Value) Value) (hadError bool) {
	if !v.peek(0).IsFloat() || !v.peek(1).IsFloat() {
		v.runtimeError("Operands must be numbers")
		return true
	}
	b := v.pop()
	a := v.pop()
	v.push(op(a, b))

	return false
}

func (v *VM) negate() (hadError bool) {
	if !v.peek(0).IsFloat() {
		v.runtimeError("Operand must be a number")
		return true
	}
	v.push(value.Float(-v.pop().AsFloat()))
	return false
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
