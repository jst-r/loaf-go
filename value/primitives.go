package value

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
