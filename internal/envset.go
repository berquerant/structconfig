package internal

import (
	"log/slog"
	"reflect"
)

// EnvReceptor sets environment variable value to the struct field.
//
// ptr should be a pointer of struct.
func EnvReceptor(
	ptr any,
	anyCallback func(StructField, string, func() reflect.Value) error,
	logger *slog.Logger,
) (*PairsReceptor, error) {
	type srcInfo struct {
		source string
		envVar string
	}
	sources := make(map[string]srcInfo)

	get := func(s StructField) (string, error) {
		name, ok := s.Tag().Name()
		if !ok {
			// ignore the field
			return "", ErrSkipParse
		}
		envVar := NewEnvVar(name)
		if v, ok := envVar.Get(); ok {
			sources[s.Name()] = srcInfo{source: "env", envVar: string(envVar)}
			return v, nil
		}
		if v, ok := s.Tag().Default(); ok {
			sources[s.Name()] = srcInfo{source: "default", envVar: ""}
			return v, nil
		}
		return "", ErrSkipParse
	}

	var onSet func(StructField, any)
	if logger != nil {
		onSet = func(s StructField, v any) {
			info := sources[s.Name()]
			attrs := []any{
				slog.String("field", s.Name()),
				slog.String("source", info.source),
				slog.Any("value", v),
			}
			if info.envVar != "" {
				attrs = append(attrs, slog.String("env_var", info.envVar))
			}
			logger.Debug("structconfig: set field", attrs...)
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
