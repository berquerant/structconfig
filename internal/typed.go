package internal

type TypedReceptorFunc[T any] func(StructField, T) error

func (f TypedReceptorFunc[T]) Call(s StructField, v T) error {
	if f == nil {
		return nil
	}
	return f(s, v)
}

type TypedReceptor interface {
	Bool(StructField, bool) error
	Int(StructField, int) error
	Int8(StructField, int8) error
	Int16(StructField, int16) error
	Int32(StructField, int32) error
	Int64(StructField, int64) error
	Uint(StructField, uint) error
	Uint8(StructField, uint8) error
	Uint16(StructField, uint16) error
	Uint32(StructField, uint32) error
	Uint64(StructField, uint64) error
	Float32(StructField, float32) error
	Float64(StructField, float64) error
	String(StructField, string) error
	Any(StructField, string) error
}

var _ TypedReceptor = &DefaultTypedReceptor{}

type DefaultTypedReceptor struct {
	BoolFunc    TypedReceptorFunc[bool]
	IntFunc     TypedReceptorFunc[int]
	Int8Func    TypedReceptorFunc[int8]
	Int16Func   TypedReceptorFunc[int16]
	Int32Func   TypedReceptorFunc[int32]
	Int64Func   TypedReceptorFunc[int64]
	UintFunc    TypedReceptorFunc[uint]
	Uint8Func   TypedReceptorFunc[uint8]
	Uint16Func  TypedReceptorFunc[uint16]
	Uint32Func  TypedReceptorFunc[uint32]
	Uint64Func  TypedReceptorFunc[uint64]
	Float32Func TypedReceptorFunc[float32]
	Float64Func TypedReceptorFunc[float64]
	StringFunc  TypedReceptorFunc[string]
	AnyFunc     TypedReceptorFunc[string]
}

func (r DefaultTypedReceptor) Bool(s StructField, v bool) error     { return r.BoolFunc.Call(s, v) }
func (r DefaultTypedReceptor) Int(s StructField, v int) error       { return r.IntFunc.Call(s, v) }
func (r DefaultTypedReceptor) Int8(s StructField, v int8) error     { return r.Int8Func.Call(s, v) }
func (r DefaultTypedReceptor) Int16(s StructField, v int16) error   { return r.Int16Func.Call(s, v) }
func (r DefaultTypedReceptor) Int32(s StructField, v int32) error   { return r.Int32Func.Call(s, v) }
func (r DefaultTypedReceptor) Int64(s StructField, v int64) error   { return r.Int64Func.Call(s, v) }
func (r DefaultTypedReceptor) Uint(s StructField, v uint) error     { return r.UintFunc.Call(s, v) }
func (r DefaultTypedReceptor) Uint8(s StructField, v uint8) error   { return r.Uint8Func.Call(s, v) }
func (r DefaultTypedReceptor) Uint16(s StructField, v uint16) error { return r.Uint16Func.Call(s, v) }
func (r DefaultTypedReceptor) Uint32(s StructField, v uint32) error { return r.Uint32Func.Call(s, v) }
func (r DefaultTypedReceptor) Uint64(s StructField, v uint64) error { return r.Uint64Func.Call(s, v) }
func (r DefaultTypedReceptor) Float32(s StructField, v float32) error {
	return r.Float32Func.Call(s, v)
}
func (r DefaultTypedReceptor) Float64(s StructField, v float64) error {
	return r.Float64Func.Call(s, v)
}
func (r DefaultTypedReceptor) String(s StructField, v string) error { return r.StringFunc.Call(s, v) }
func (r DefaultTypedReceptor) Any(s StructField, v string) error    { return r.AnyFunc.Call(s, v) }
