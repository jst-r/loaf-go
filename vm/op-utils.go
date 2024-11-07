package vm

import "github.com/jst-r/loaf-go/value"

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
