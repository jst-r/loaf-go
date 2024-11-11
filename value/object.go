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

func object[T any](val *T) Value {
	return Value{ValueTypeObject, unsafeBitCast[*T, uint64](val)}
}

func String(val string) Value {
	return object(&ObjString{Obj{ObjTypeString}, val})
}

func (v Value) IsObject() bool {
	return v.t == ValueTypeObject
}

func (v Value) IsObjectType(t ObjType) bool {
	return v.IsObject() && v.AsObject().t == t
}

// Unsafe, caller must verify the type before calling this method
func (v Value) AsObject() *Obj {
	return unsafeBitCast[uint64, *Obj](v.mem)
}

// Unsafe, caller must verify the type before calling this method
func (v Value) AsString() *ObjString {
	return unsafeBitCast[uint64, *ObjString](v.mem)
}
