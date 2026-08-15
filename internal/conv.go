package internal

import (
	"strconv"
	"time"
)

// Converter is a set of string conversions.
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
	BoolSlice(s string) ([]bool, error)
	IntSlice(s string) ([]int, error)
	Int32Slice(s string) ([]int32, error)
	Int64Slice(s string) ([]int64, error)
	UintSlice(s string) ([]uint, error)
	Float32Slice(s string) ([]float32, error)
	Float64Slice(s string) ([]float64, error)
	StringSlice(s string) ([]string, error)
	Duration(s string) (time.Duration, error)
	Time(s string) (time.Time, error)
	Count(s string) (int, error)
}

// ConvFunc is a function that converts string into T.
type ConvFunc[T any] func(string) (T, error)

// Call calls self if not nil, otherwise returns the default value.
func (f ConvFunc[T]) Call(s string) (T, error) {
	if f == nil {
		var t T
		return t, nil
	}
	return f(s)
}

var _ Converter = &DefaultConverter{}

// DefaultConverter implements [Converter].
type DefaultConverter struct {
	BoolFunc         ConvFunc[bool]
	IntFunc          ConvFunc[int]
	Int8Func         ConvFunc[int8]
	Int16Func        ConvFunc[int16]
	Int32Func        ConvFunc[int32]
	Int64Func        ConvFunc[int64]
	UintFunc         ConvFunc[uint]
	Uint8Func        ConvFunc[uint8]
	Uint16Func       ConvFunc[uint16]
	Uint32Func       ConvFunc[uint32]
	Uint64Func       ConvFunc[uint64]
	Float32Func      ConvFunc[float32]
	Float64Func      ConvFunc[float64]
	StringFunc       ConvFunc[string]
	BoolSliceFunc    ConvFunc[[]bool]
	IntSliceFunc     ConvFunc[[]int]
	Int32SliceFunc   ConvFunc[[]int32]
	Int64SliceFunc   ConvFunc[[]int64]
	UintSliceFunc    ConvFunc[[]uint]
	Float32SliceFunc ConvFunc[[]float32]
	Float64SliceFunc ConvFunc[[]float64]
	StringSliceFunc  ConvFunc[[]string]
	DurationFunc     ConvFunc[time.Duration]
	TimeFunc         ConvFunc[time.Time]
	CountFunc        ConvFunc[int]
}

func (c DefaultConverter) Bool(s string) (bool, error)                  { return c.BoolFunc.Call(s) }
func (c DefaultConverter) Int(s string) (int, error)                    { return c.IntFunc.Call(s) }
func (c DefaultConverter) Int8(s string) (int8, error)                  { return c.Int8Func.Call(s) }
func (c DefaultConverter) Int16(s string) (int16, error)                { return c.Int16Func.Call(s) }
func (c DefaultConverter) Int32(s string) (int32, error)                { return c.Int32Func.Call(s) }
func (c DefaultConverter) Int64(s string) (int64, error)                { return c.Int64Func.Call(s) }
func (c DefaultConverter) Uint(s string) (uint, error)                  { return c.UintFunc.Call(s) }
func (c DefaultConverter) Uint8(s string) (uint8, error)                { return c.Uint8Func.Call(s) }
func (c DefaultConverter) Uint16(s string) (uint16, error)              { return c.Uint16Func.Call(s) }
func (c DefaultConverter) Uint32(s string) (uint32, error)              { return c.Uint32Func.Call(s) }
func (c DefaultConverter) Uint64(s string) (uint64, error)              { return c.Uint64Func.Call(s) }
func (c DefaultConverter) Float32(s string) (float32, error)            { return c.Float32Func.Call(s) }
func (c DefaultConverter) Float64(s string) (float64, error)            { return c.Float64Func.Call(s) }
func (c DefaultConverter) String(s string) (string, error)              { return c.StringFunc.Call(s) }
func (c DefaultConverter) BoolSlice(s string) ([]bool, error)          { return c.BoolSliceFunc.Call(s) }
func (c DefaultConverter) IntSlice(s string) ([]int, error)            { return c.IntSliceFunc.Call(s) }
func (c DefaultConverter) Int32Slice(s string) ([]int32, error)        { return c.Int32SliceFunc.Call(s) }
func (c DefaultConverter) Int64Slice(s string) ([]int64, error)        { return c.Int64SliceFunc.Call(s) }
func (c DefaultConverter) UintSlice(s string) ([]uint, error)          { return c.UintSliceFunc.Call(s) }
func (c DefaultConverter) Float32Slice(s string) ([]float32, error)    { return c.Float32SliceFunc.Call(s) }
func (c DefaultConverter) Float64Slice(s string) ([]float64, error)    { return c.Float64SliceFunc.Call(s) }
func (c DefaultConverter) StringSlice(s string) ([]string, error)      { return c.StringSliceFunc.Call(s) }
func (c DefaultConverter) Duration(s string) (time.Duration, error)    { return c.DurationFunc.Call(s) }
func (c DefaultConverter) Time(s string) (time.Time, error)            { return c.TimeFunc.Call(s) }
func (c DefaultConverter) Count(s string) (int, error)                  { return c.CountFunc.Call(s) }

func NewConv() *DefaultConverter {
	return &DefaultConverter{
		BoolFunc:         strconv.ParseBool,
		IntFunc:          ParseInt[int],
		Int8Func:         ParseInt[int8],
		Int16Func:        ParseInt[int16],
		Int32Func:        ParseInt[int32],
		Int64Func:        ParseInt[int64],
		UintFunc:         ParseUint[uint],
		Uint8Func:        ParseUint[uint8],
		Uint16Func:       ParseUint[uint16],
		Uint32Func:       ParseUint[uint32],
		Uint64Func:       ParseUint[uint64],
		Float32Func:      ParseFloat[float32],
		Float64Func:      ParseFloat[float64],
		StringFunc:       func(s string) (string, error) { return s, nil },
		BoolSliceFunc:    func(s string) ([]bool, error) { return ParseSlice(s, strconv.ParseBool) },
		IntSliceFunc:     func(s string) ([]int, error) { return ParseSlice(s, ParseInt[int]) },
		Int32SliceFunc:   func(s string) ([]int32, error) { return ParseSlice(s, ParseInt[int32]) },
		Int64SliceFunc:   func(s string) ([]int64, error) { return ParseSlice(s, ParseInt[int64]) },
		UintSliceFunc:    func(s string) ([]uint, error) { return ParseSlice(s, ParseUint[uint]) },
		Float32SliceFunc: func(s string) ([]float32, error) { return ParseSlice(s, ParseFloat[float32]) },
		Float64SliceFunc: func(s string) ([]float64, error) { return ParseSlice(s, ParseFloat[float64]) },
		StringSliceFunc:  func(s string) ([]string, error) { return ParseSlice(s, func(x string) (string, error) { return x, nil }) },
		DurationFunc:     ParseDuration,
		TimeFunc:         ParseTime,
		CountFunc:        ParseInt[int],
	}
}
