package value

import "unsafe"

type ObjPool struct {
	objects []unsafe.Pointer
}

func NewObjPool() *ObjPool {
	return &ObjPool{objects: make([]unsafe.Pointer, 0)}
}

func (p *ObjPool) NewString(val string) Value {
	obj := &ObjString{objMetadata{ObjTypeString}, val}
	p.objects = append(p.objects, unsafe.Pointer(obj))
	addr := uint64(unsafeBitCast[*ObjString, uintptr](obj))
	return Value{ValueTypeObject, addr}
}

func (p *ObjPool) Append(other *ObjPool) {
	p.objects = append(p.objects, other.objects...)
}
