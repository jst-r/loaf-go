package value

type objMetadata struct {
	t ObjType
}

type ObjString struct {
	objMetadata
	Str string
}

type ObjType int

const (
	ObjTypeString ObjType = iota
)

func (v Value) IsObject() bool {
	return v.t == ValueTypeObject
}

func (v Value) IsObjectType(t ObjType) bool {
	return v.IsObject() && v.ObjectType() == t
}

func (v Value) IsString() bool {
	return v.IsObjectType(ObjTypeString)
}

// Unsafe, caller must verify the type before calling this method
func (v Value) ObjectType() ObjType {
	return v.AsObject().t
}

// Unsafe, caller must verify the type before calling this method
func (v Value) AsObject() *objMetadata {
	return unsafeBitCast[uint64, *objMetadata](v.mem)
}

// Unsafe, caller must verify the type before calling this method
func (v Value) AsString() *ObjString {
	return unsafeBitCast[uint64, *ObjString](v.mem)
}
