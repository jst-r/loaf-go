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
	ValueTypeObject
)

type Value struct {
	t   ValueType
	mem uint64 // 8 bytes of contiguous memory
}

type ValueArray []Value

func (v Value) IsFalsey() bool {
	return v.IsNil() || (v.IsBool() && !v.AsBool())
}

func ValuesEqual(a, b Value) bool {
	if a.t != b.t {
		return false
	}
	if !a.IsObject() {
		return a.mem == b.mem // Since go initializes all memory to zero this works for smaller types (e.g. bool)
	}

	switch a.ObjectType() {
	case ObjTypeString:
		return a.AsString().Str == b.AsString().Str
	default:
		return false
	}
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
		return fmt.Sprintf("%v", v.AsFloat())
	case ValueTypeObject:
		switch v.ObjectType() {
		case ObjTypeString:
			return v.AsString().Str
		default:
			panic("unreachable")
		}
	default:
		panic("unreachable")
	}
}

func unsafeBitCast[A, B any](a A) B {
	return *(*B)(unsafe.Pointer(&a))
}
