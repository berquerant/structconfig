package internal

import (
	"log/slog"
	"reflect"

	"github.com/spf13/pflag"
)

// PFlagSetReceptor returns a [Receptor] that can define the command-line flags.
func PFlagSetReceptor(fs *pflag.FlagSet) *PairsReceptor {
	return FlagSetReceptor(PFlagSetTypeReceptor(fs))
}

// PFlagGetReceptor returns a [Receptor] that can retrieve values from the parsed command-line flags.
func PFlagGetReceptor(
	ptr any,
	fs *pflag.FlagSet,
	anyCallback func(StructField, string, func() reflect.Value) error,
	logger *slog.Logger,
) (*PairsReceptor, error) {
	type flagInfo struct {
		flagName string
		changed  bool
	}
	sources := make(map[string]flagInfo)

	get := func(s StructField) (string, error) {
		if name, ok := s.Tag().Name(); ok {
			sources[s.Name()] = flagInfo{
				flagName: name,
				changed:  fs.Changed(name),
			}
			return name, nil
		}
		return "", ErrParseAsDefault
	}

	var onSet func(StructField, any)
	if logger != nil {
		onSet = func(s StructField, v any) {
			info := sources[s.Name()]
			source := "flag"
			if !info.changed {
				source = "flag_default"
			}
			logger.Debug(
				"structconfig: set field",
				slog.String("field", s.Name()),
				slog.String("source", source),
				slog.String("flag", info.flagName),
				slog.Any("value", v),
				slog.Bool("changed", info.changed),
			)
		}
	}

	typedReceptor, err := SetTypedReceptor(ptr, anyCallback, onSet)
	if err != nil {
		return nil, err
	}
	return PairsSynthReceptor(
		get,
		PFlagGetConverter(fs),
		typedReceptor,
	), nil
}

func pflagSetFunc[T any](
	f func(string, T, string) *T,
	g func(string, string, T, string) *T,
) TypedReceptorFunc[T] {
	return func(s StructField, defaultValue T) error {
		if name, ok := s.Tag().Name(); ok {
			if short, ok := s.Tag().Short(); ok {
				_ = g(name, short, defaultValue, s.Tag().Usage())
				return nil
			}
			_ = f(name, defaultValue, s.Tag().Usage())
		}
		return nil
	}
}

func PFlagSetTypeReceptor(fs *pflag.FlagSet) *DefaultTypedReceptor {
	return &DefaultTypedReceptor{
		BoolFunc:    pflagSetFunc(fs.Bool, fs.BoolP),
		IntFunc:     pflagSetFunc(fs.Int, fs.IntP),
		Int8Func:    pflagSetFunc(fs.Int8, fs.Int8P),
		Int16Func:   pflagSetFunc(fs.Int16, fs.Int16P),
		Int32Func:   pflagSetFunc(fs.Int32, fs.Int32P),
		Int64Func:   pflagSetFunc(fs.Int64, fs.Int64P),
		UintFunc:    pflagSetFunc(fs.Uint, fs.UintP),
		Uint8Func:   pflagSetFunc(fs.Uint8, fs.Uint8P),
		Uint16Func:  pflagSetFunc(fs.Uint16, fs.Uint16P),
		Uint32Func:  pflagSetFunc(fs.Uint32, fs.Uint32P),
		Uint64Func:  pflagSetFunc(fs.Uint64, fs.Uint64P),
		Float32Func: pflagSetFunc(fs.Float32, fs.Float32P),
		Float64Func: pflagSetFunc(fs.Float64, fs.Float64P),
		StringFunc:  pflagSetFunc(fs.String, fs.StringP),
		AnyFunc:     pflagSetFunc(fs.String, fs.StringP),
	}
}

func PFlagGetConverter(fs *pflag.FlagSet) *DefaultConverter {
	return &DefaultConverter{
		BoolFunc:    fs.GetBool,
		IntFunc:     fs.GetInt,
		Int8Func:    fs.GetInt8,
		Int16Func:   fs.GetInt16,
		Int32Func:   fs.GetInt32,
		Int64Func:   fs.GetInt64,
		UintFunc:    fs.GetUint,
		Uint8Func:   fs.GetUint8,
		Uint16Func:  fs.GetUint16,
		Uint32Func:  fs.GetUint32,
		Uint64Func:  fs.GetUint64,
		Float32Func: fs.GetFloat32,
		Float64Func: fs.GetFloat64,
		StringFunc:  fs.GetString,
	}
}
