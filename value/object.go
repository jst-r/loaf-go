package value

type Obj struct {
	t ObjType
}

type ObjString struct {
	Obj
	Str string
}

type ObjType int

const (
	ObjTypeString ObjType = iota
)

// Unsafe, caller must verify the type before calling this method
func Object(val *interface{}) Value {
	return Value{ValueTypeObject, unsafeBitCast[*interface{}, uint64](val)}
}

func (v Value) IsObject() bool {
	return v.t == ValueTypeObject
}

func (v Value) IsObjectType(t ObjType) bool {
	return v.IsObject() && v.AsObject().t == t
}

func (v Value) AsObject() *Obj {
	return unsafeBitCast[uint64, *Obj](v.mem)
}

func (v Value) AsString() *ObjString {
	return unsafeBitCast[uint64, *ObjString](v.mem)
}
