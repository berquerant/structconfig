package internal

import "reflect"

// Call calls the appropriate method of r for the kind of f.
func Call(r Receptor, f StructField) error {
	return Switch(r, f.Kind())(f)
}

// Switch chooses a method of r that corresponds to kind.
func Switch(r Receptor, kind reflect.Kind) func(StructField) error {
	switch kind {
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
	default:
		return r.Any
	}
}
