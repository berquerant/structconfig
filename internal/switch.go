package internal

import (
	"reflect"
	"time"
)

var (
	durationTyp = reflect.TypeOf(time.Duration(0))
	timeTyp     = reflect.TypeOf(time.Time{})
)

// Call calls the appropriate method of r for the field f.
func Call(r Receptor, f StructField) error {
	return Switch(r, f)(f)
}

// Switch chooses a method of r that corresponds to f.
func Switch(r Receptor, f StructField) func(StructField) error {
	if f.Tag() != nil && f.Tag().Count() {
		return r.Count
	}

	t := f.RType()
	if t == durationTyp {
		return r.Duration
	}
	if t == timeTyp {
		return r.Time
	}

	switch f.Kind() {
	case reflect.Bool:
		return r.Bool
	case reflect.Int:
		return r.Int
	case reflect.Int8:
		return r.Int8
	case reflect.Int16:
		return r.Int16
	case reflect.Int32:
		return r.Int32
	case reflect.Int64:
		return r.Int64
	case reflect.Uint:
		return r.Uint
	case reflect.Uint8:
		return r.Uint8
	case reflect.Uint16:
		return r.Uint16
	case reflect.Uint32:
		return r.Uint32
	case reflect.Uint64:
		return r.Uint64
	case reflect.Float32:
		return r.Float32
	case reflect.Float64:
		return r.Float64
	case reflect.String:
		return r.String
	case reflect.Slice:
		if t != nil {
			switch t.Elem().Kind() {
			case reflect.Bool:
				return r.BoolSlice
			case reflect.Int:
				return r.IntSlice
			case reflect.Int32:
				return r.Int32Slice
			case reflect.Int64:
				return r.Int64Slice
			case reflect.Uint:
				return r.UintSlice
			case reflect.Float32:
				return r.Float32Slice
			case reflect.Float64:
				return r.Float64Slice
			case reflect.String:
				return r.StringSlice
			}
		}
		return r.Any
	default:
		return r.Any
	}
}
