package internal

import "reflect"

// DefaultReceptor sets default tag value to the struct field.
//
// ptr should be a pointer of struct.
func DefaultReceptor(
	ptr any,
	anyCallback func(StructField, string, func() reflect.Value) error,
) (*PairsReceptor, error) {
	get := func(s StructField) (string, error) {
		if v, ok := s.Tag().Default(); ok {
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
