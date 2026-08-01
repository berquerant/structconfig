package internal

import (
	"log/slog"
	"reflect"
)

// DefaultReceptor sets default tag value to the struct field.
//
// ptr should be a pointer of struct.
func DefaultReceptor(
	ptr any,
	anyCallback func(StructField, string, func() reflect.Value) error,
	logger *slog.Logger,
) (*PairsReceptor, error) {
	get := func(s StructField) (string, error) {
		if v, ok := s.Tag().Default(); ok {
			return v, nil
		}
		return "", ErrSkipParse
	}

	var onSet func(StructField, any)
	if logger != nil {
		onSet = func(s StructField, v any) {
			logger.Debug(
				"structconfig: set field",
				slog.String("field", s.Name()),
				slog.String("source", "default"),
				slog.Any("value", v),
			)
		}
	}

	return SetReceptor(
		ptr,
		get,
		NewConv(),
		anyCallback,
		onSet,
	)
}
