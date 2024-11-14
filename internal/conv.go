package internal

import "strconv"

type Converter interface {
	Bool(s string) (bool, error)
	Int(s string) (int, error)
	Int8(s string) (int8, error)
	Int16(s string) (int16, error)
	Int32(s string) (int32, error)
	Int64(s string) (int64, error)
	Uint(s string) (uint, error)
	Uint8(s string) (uint8, error)
	Uint16(s string) (uint16, error)
	Uint32(s string) (uint32, error)
	Uint64(s string) (uint64, error)
	Float32(s string) (float32, error)
	Float64(s string) (float64, error)
	String(s string) (string, error)
}

type ConvFunc[T any] func(string) (T, error)

func (f ConvFunc[T]) Call(s string) (T, error) {
	if f == nil {
		var t T
		return t, nil
	}
	return f(s)
}

var _ Converter = &DefaultConverter{}

type DefaultConverter struct {
	BoolFunc    ConvFunc[bool]
	IntFunc     ConvFunc[int]
	Int8Func    ConvFunc[int8]
	Int16Func   ConvFunc[int16]
	Int32Func   ConvFunc[int32]
	Int64Func   ConvFunc[int64]
	UintFunc    ConvFunc[uint]
	Uint8Func   ConvFunc[uint8]
	Uint16Func  ConvFunc[uint16]
	Uint32Func  ConvFunc[uint32]
	Uint64Func  ConvFunc[uint64]
	Float32Func ConvFunc[float32]
	Float64Func ConvFunc[float64]
	StringFunc  ConvFunc[string]
}

func (c DefaultConverter) Bool(s string) (bool, error)       { return c.BoolFunc.Call(s) }
func (c DefaultConverter) Int(s string) (int, error)         { return c.IntFunc.Call(s) }
func (c DefaultConverter) Int8(s string) (int8, error)       { return c.Int8Func.Call(s) }
func (c DefaultConverter) Int16(s string) (int16, error)     { return c.Int16Func.Call(s) }
func (c DefaultConverter) Int32(s string) (int32, error)     { return c.Int32Func.Call(s) }
func (c DefaultConverter) Int64(s string) (int64, error)     { return c.Int64Func.Call(s) }
func (c DefaultConverter) Uint(s string) (uint, error)       { return c.UintFunc.Call(s) }
func (c DefaultConverter) Uint8(s string) (uint8, error)     { return c.Uint8Func.Call(s) }
func (c DefaultConverter) Uint16(s string) (uint16, error)   { return c.Uint16Func.Call(s) }
func (c DefaultConverter) Uint32(s string) (uint32, error)   { return c.Uint32Func.Call(s) }
func (c DefaultConverter) Uint64(s string) (uint64, error)   { return c.Uint64Func.Call(s) }
func (c DefaultConverter) Float32(s string) (float32, error) { return c.Float32Func.Call(s) }
func (c DefaultConverter) Float64(s string) (float64, error) { return c.Float64Func.Call(s) }
func (c DefaultConverter) String(s string) (string, error)   { return c.StringFunc.Call(s) }

func NewConv() *DefaultConverter {
	return &DefaultConverter{
		BoolFunc:    strconv.ParseBool,
		IntFunc:     ParseInt[int],
		Int8Func:    ParseInt[int8],
		Int16Func:   ParseInt[int16],
		Int32Func:   ParseInt[int32],
		Int64Func:   ParseInt[int64],
		UintFunc:    ParseUint[uint],
		Uint8Func:   ParseUint[uint8],
		Uint16Func:  ParseUint[uint16],
		Uint32Func:  ParseUint[uint32],
		Uint64Func:  ParseUint[uint64],
		Float32Func: ParseFloat[float32],
		Float64Func: ParseFloat[float64],
		StringFunc:  func(s string) (string, error) { return s, nil },
	}
}
