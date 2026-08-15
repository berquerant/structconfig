package internal

//go:generate go tool dataclass -type StructField -field "Name string|Kind reflect.Kind|Tag *Tag|RType reflect.Type" -output structfield_dataclass_generated.go

// Receptor accepts [StructField].
type Receptor interface {
	BoolReceptor
	IntReceptor
	UintReceptor
	FloatReceptor
	StringReceptor
	SliceReceptor
	DurationReceptor
	TimeReceptor
	CountReceptor
	AnyReceptor
}

type BoolReceptor interface {
	Bool(StructField) error
}

type IntReceptor interface {
	Int(StructField) error
	Int8(StructField) error
	Int16(StructField) error
	Int32(StructField) error
	Int64(StructField) error
}

type UintReceptor interface {
	Uint(StructField) error
	Uint8(StructField) error
	Uint16(StructField) error
	Uint32(StructField) error
	Uint64(StructField) error
}

type FloatReceptor interface {
	Float32(StructField) error
	Float64(StructField) error
}

type StringReceptor interface {
	String(StructField) error
}

type SliceReceptor interface {
	BoolSlice(StructField) error
	IntSlice(StructField) error
	Int32Slice(StructField) error
	Int64Slice(StructField) error
	UintSlice(StructField) error
	Float32Slice(StructField) error
	Float64Slice(StructField) error
	StringSlice(StructField) error
}

type DurationReceptor interface {
	Duration(StructField) error
}

type TimeReceptor interface {
	Time(StructField) error
}

type CountReceptor interface {
	Count(StructField) error
}

type AnyReceptor interface {
	Any(StructField) error
}
