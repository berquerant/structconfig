package internal

import "reflect"

// EnvReceptor sets environment variable value to the struct field.
//
// ptr should be a pointer of struct.
func EnvReceptor(
	ptr any,
	anyCallback func(StructField, string, func() reflect.Value) error,
) (*PairsReceptor, error) {
	get := func(s StructField) (string, error) {
		name, ok := s.Tag().Name()
		if !ok {
			// ignore the field
			return "", ErrSkipParse
		}
		if v, ok := NewEnvVar(name).Get(); ok {
			return v, nil
		}
		return "", ErrSkipParse
	}

	return SetReceptor(
		ptr,
		get,
		NewConv(),
		anyCallback,
	)
}
