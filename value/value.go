package value

import (
	"fmt"
	"unsafe"
)

type ValueType int

const (
	ValueTypeNil ValueType = iota
	ValueTypeBool
	ValueTypeFloat
)

type Value struct {
	t   ValueType
	mem uint64 // 8 bytes of contiguous memory
}

type ValueArray []Value

var Nil = Value{ValueTypeNil, 0}

func Float(val float64) Value {
	return Value{ValueTypeFloat, unsafeBitCast[float64, uint64](val)}
}

// Unsafe, caller must verify the type before calling this method
func (v Value) AsFloat() float64 {
	return unsafeBitCast[uint64, float64](v.mem)
}

func Bool(val bool) Value {
	if val {
		return Value{ValueTypeBool, 1}
	} else {
		return Value{ValueTypeBool, 0}
	}
}

func (v Value) AsBool() bool {
	return v.mem != 0
}

func (v Value) IsFalsey() bool {
	return v.IsNil() || (v.IsBool() && !v.AsBool())
}

func (v Value) FormatString() string {
	switch v.t {
	case ValueTypeNil:
		return "nil"
	case ValueTypeBool:
		if v.AsBool() {
			return "true"
		} else {
			return "false"
		}
	case ValueTypeFloat:
		return fmt.Sprintf("%f", v.AsFloat())
	default:
		panic("unreachable")
	}
}

func (v Value) Is(t ValueType) bool {
	return v.t == t
}

func (v Value) IsNil() bool {
	return v.t == ValueTypeNil
}

func (v Value) IsBool() bool {
	return v.t == ValueTypeBool
}

func (v Value) IsFloat() bool {
	return v.t == ValueTypeFloat
}

func unsafeBitCast[A, B any](a A) B {
	return *(*B)(unsafe.Pointer(&a))
}
