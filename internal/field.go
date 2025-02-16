package internal

//go:generate go tool dataclass -type StructField -field "Name string|Kind reflect.Kind|Tag *Tag" -output structfield_dataclass_generated.go

// Receptor accepts [StructField].
type Receptor interface {
	BoolReceptor
	IntReceptor
	UintReceptor
	FloatReceptor
	StringReceptor
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

type AnyReceptor interface {
	Any(StructField) error
}
