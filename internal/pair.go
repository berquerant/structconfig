package internal

import (
	"errors"
	"time"
)

var (
	// Result of conv as default value in ParsePair.
	ErrParseAsDefault = errors.New("ParseAsDefault")
)

// NewPairSynth constructs a pair of functions to process [StructField].
//   - get: extract string from [StructField]
//   - conv: convert the value
//   - callback: accept the converted value
//
// get and conv can return [ErrParseAsDefault] or [ErrSkipParse].
func NewPairSynth[T any](
	get func(StructField) (string, error),
	conv func(string) (T, error),
	callback func(StructField, T) error,
) *ParsePair[T] {
	return &ParsePair[T]{
		conv: func(s StructField) (T, error) {
			x, err := get(s)
			switch {
			case err == nil:
				return conv(x)
			case errors.Is(err, ErrParseAsDefault):
				var t T
				return t, nil
			default:
				var t T
				return t, err
			}
		},
		callback: callback,
	}
}

func NewParsePair[T any](
	conv func(StructField) (T, error),
	callback func(StructField, T) error,
) *ParsePair[T] {
	return &ParsePair[T]{
		conv:     conv,
		callback: callback,
	}
}

// ParsePair defines a conversion of [StructField] that can fail.
type ParsePair[T any] struct {
	conv     func(StructField) (T, error)
	callback func(StructField, T) error
}

func (p *ParsePair[T]) Try(s StructField) error {
	if p == nil {
		return nil
	}
	return TryParse(p.conv, p.callback)(s)
}

var _ Receptor = &PairsReceptor{}

// PairsReceptor is a set of [ParsePair], implements [Receptor].
type PairsReceptor struct {
	BoolPair         *ParsePair[bool]
	IntPair          *ParsePair[int]
	Int8Pair         *ParsePair[int8]
	Int16Pair        *ParsePair[int16]
	Int32Pair        *ParsePair[int32]
	Int64Pair        *ParsePair[int64]
	UintPair         *ParsePair[uint]
	Uint8Pair        *ParsePair[uint8]
	Uint16Pair       *ParsePair[uint16]
	Uint32Pair       *ParsePair[uint32]
	Uint64Pair       *ParsePair[uint64]
	Float32Pair      *ParsePair[float32]
	Float64Pair      *ParsePair[float64]
	StringPair       *ParsePair[string]
	BoolSlicePair    *ParsePair[[]bool]
	IntSlicePair     *ParsePair[[]int]
	Int32SlicePair   *ParsePair[[]int32]
	Int64SlicePair   *ParsePair[[]int64]
	UintSlicePair    *ParsePair[[]uint]
	Float32SlicePair *ParsePair[[]float32]
	Float64SlicePair *ParsePair[[]float64]
	StringSlicePair  *ParsePair[[]string]
	DurationPair     *ParsePair[time.Duration]
	TimePair         *ParsePair[time.Time]
	CountPair        *ParsePair[int]
	AnyPair          *ParsePair[string]
}

func (r PairsReceptor) Bool(f StructField) error         { return r.BoolPair.Try(f) }
func (r PairsReceptor) Int(f StructField) error          { return r.IntPair.Try(f) }
func (r PairsReceptor) Int8(f StructField) error         { return r.Int8Pair.Try(f) }
func (r PairsReceptor) Int16(f StructField) error        { return r.Int16Pair.Try(f) }
func (r PairsReceptor) Int32(f StructField) error        { return r.Int32Pair.Try(f) }
func (r PairsReceptor) Int64(f StructField) error        { return r.Int64Pair.Try(f) }
func (r PairsReceptor) Uint(f StructField) error         { return r.UintPair.Try(f) }
func (r PairsReceptor) Uint8(f StructField) error        { return r.Uint8Pair.Try(f) }
func (r PairsReceptor) Uint16(f StructField) error       { return r.Uint16Pair.Try(f) }
func (r PairsReceptor) Uint32(f StructField) error       { return r.Uint32Pair.Try(f) }
func (r PairsReceptor) Uint64(f StructField) error       { return r.Uint64Pair.Try(f) }
func (r PairsReceptor) Float32(f StructField) error      { return r.Float32Pair.Try(f) }
func (r PairsReceptor) Float64(f StructField) error      { return r.Float64Pair.Try(f) }
func (r PairsReceptor) String(f StructField) error       { return r.StringPair.Try(f) }
func (r PairsReceptor) BoolSlice(f StructField) error    { return r.BoolSlicePair.Try(f) }
func (r PairsReceptor) IntSlice(f StructField) error     { return r.IntSlicePair.Try(f) }
func (r PairsReceptor) Int32Slice(f StructField) error   { return r.Int32SlicePair.Try(f) }
func (r PairsReceptor) Int64Slice(f StructField) error   { return r.Int64SlicePair.Try(f) }
func (r PairsReceptor) UintSlice(f StructField) error    { return r.UintSlicePair.Try(f) }
func (r PairsReceptor) Float32Slice(f StructField) error { return r.Float32SlicePair.Try(f) }
func (r PairsReceptor) Float64Slice(f StructField) error { return r.Float64SlicePair.Try(f) }
func (r PairsReceptor) StringSlice(f StructField) error  { return r.StringSlicePair.Try(f) }
func (r PairsReceptor) Duration(f StructField) error     { return r.DurationPair.Try(f) }
func (r PairsReceptor) Time(f StructField) error         { return r.TimePair.Try(f) }
func (r PairsReceptor) Count(f StructField) error        { return r.CountPair.Try(f) }
func (r PairsReceptor) Any(f StructField) error          { return r.AnyPair.Try(f) }

// PairsSynthReceptor synthesizes [Converter] and [TypedReceptor].
// get extracts the value from [StructField], converter converts it and typedReceptor accepts it.
func PairsSynthReceptor(
	get func(StructField) (string, error),
	converter Converter,
	typedReceptor TypedReceptor,
) *PairsReceptor {
	c := converter
	t := typedReceptor

	return &PairsReceptor{
		BoolPair:         NewPairSynth(get, c.Bool, t.Bool),
		IntPair:          NewPairSynth(get, c.Int, t.Int),
		Int8Pair:         NewPairSynth(get, c.Int8, t.Int8),
		Int16Pair:        NewPairSynth(get, c.Int16, t.Int16),
		Int32Pair:        NewPairSynth(get, c.Int32, t.Int32),
		Int64Pair:        NewPairSynth(get, c.Int64, t.Int64),
		UintPair:         NewPairSynth(get, c.Uint, t.Uint),
		Uint8Pair:        NewPairSynth(get, c.Uint8, t.Uint8),
		Uint16Pair:       NewPairSynth(get, c.Uint16, t.Uint16),
		Uint32Pair:       NewPairSynth(get, c.Uint32, t.Uint32),
		Uint64Pair:       NewPairSynth(get, c.Uint64, t.Uint64),
		Float32Pair:      NewPairSynth(get, c.Float32, t.Float32),
		Float64Pair:      NewPairSynth(get, c.Float64, t.Float64),
		StringPair:       NewPairSynth(get, c.String, t.String),
		BoolSlicePair:    NewPairSynth(get, c.BoolSlice, t.BoolSlice),
		IntSlicePair:     NewPairSynth(get, c.IntSlice, t.IntSlice),
		Int32SlicePair:   NewPairSynth(get, c.Int32Slice, t.Int32Slice),
		Int64SlicePair:   NewPairSynth(get, c.Int64Slice, t.Int64Slice),
		UintSlicePair:    NewPairSynth(get, c.UintSlice, t.UintSlice),
		Float32SlicePair: NewPairSynth(get, c.Float32Slice, t.Float32Slice),
		Float64SlicePair: NewPairSynth(get, c.Float64Slice, t.Float64Slice),
		StringSlicePair:  NewPairSynth(get, c.StringSlice, t.StringSlice),
		DurationPair:     NewPairSynth(get, c.Duration, t.Duration),
		TimePair:         NewPairSynth(get, c.Time, t.Time),
		CountPair:        NewPairSynth(get, c.Count, t.Count),
		AnyPair:          NewPairSynth(get, c.String, t.Any),
	}
}
