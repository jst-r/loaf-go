package value

import "fmt"

type Value struct {
	val float64
}

type ValueArray []Value

func Float(val float64) Value {
	return Value{val}
}

func (v Value) AsFloat() float64 {
	return v.val
}

func (v Value) String() string {
	return fmt.Sprintf("%f", v.val)
}
