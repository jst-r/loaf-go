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

func (v *VM) add() {
	if v.peek(0).IsFloat() || v.peek(1).IsFloat() {
		b := v.pop()
		a := v.pop()
		v.push(value.Float(a.AsFloat() + b.AsFloat()))
	} else if v.peek(0).IsString() || v.peek(1).IsString() {
		b := v.pop()
		a := v.pop()
		v.push(value.String(a.AsString().Str + b.AsString().Str))
	} else {
		panic("Type error: operands must be numbers or strings")
	}
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
