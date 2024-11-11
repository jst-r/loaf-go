package vm

import "github.com/jst-r/loaf-go/value"

func (v *VM) binaryOp(op func(a, b Value) Value) {
	if !v.peek(0).IsFloat() || !v.peek(1).IsFloat() {
		panic("Type error: operands must be numbers")
	}
	b := v.pop()
	a := v.pop()
	v.push(op(a, b))
}

func (v *VM) negate() {
	if !v.peek(0).IsFloat() {
		panic("Type error: operand must be a number")
	}
	v.push(value.Float(-v.pop().AsFloat()))
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

func greater(a, b Value) Value {
	return value.Bool(a.AsFloat() > b.AsFloat())
}

func less(a, b Value) Value {
	return value.Bool(a.AsFloat() < b.AsFloat())
}
