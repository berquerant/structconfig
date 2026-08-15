package internal

import (
	"log/slog"
	"reflect"
	"time"

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

func pflagCountSetFunc(
	f func(string, string, string) *int,
) TypedReceptorFunc[int] {
	return func(s StructField, defaultValue int) error {
		if name, ok := s.Tag().Name(); ok {
			short, _ := s.Tag().Short()
			_ = f(name, short, s.Tag().Usage())
		}
		return nil
	}
}

func pflagTimeSetFunc(
	f func(string, string, string) *string,
	g func(string, string, string, string) *string,
) TypedReceptorFunc[time.Time] {
	return func(s StructField, defaultValue time.Time) error {
		if name, ok := s.Tag().Name(); ok {
			defStr := ""
			if !defaultValue.IsZero() {
				defStr = defaultValue.Format(time.RFC3339)
			}
			if short, ok := s.Tag().Short(); ok {
				_ = g(name, short, defStr, s.Tag().Usage())
				return nil
			}
			_ = f(name, defStr, s.Tag().Usage())
		}
		return nil
	}
}

func PFlagSetTypeReceptor(fs *pflag.FlagSet) *DefaultTypedReceptor {
	return &DefaultTypedReceptor{
		BoolFunc:         pflagSetFunc(fs.Bool, fs.BoolP),
		IntFunc:          pflagSetFunc(fs.Int, fs.IntP),
		Int8Func:         pflagSetFunc(fs.Int8, fs.Int8P),
		Int16Func:        pflagSetFunc(fs.Int16, fs.Int16P),
		Int32Func:        pflagSetFunc(fs.Int32, fs.Int32P),
		Int64Func:        pflagSetFunc(fs.Int64, fs.Int64P),
		UintFunc:         pflagSetFunc(fs.Uint, fs.UintP),
		Uint8Func:        pflagSetFunc(fs.Uint8, fs.Uint8P),
		Uint16Func:       pflagSetFunc(fs.Uint16, fs.Uint16P),
		Uint32Func:       pflagSetFunc(fs.Uint32, fs.Uint32P),
		Uint64Func:       pflagSetFunc(fs.Uint64, fs.Uint64P),
		Float32Func:      pflagSetFunc(fs.Float32, fs.Float32P),
		Float64Func:      pflagSetFunc(fs.Float64, fs.Float64P),
		StringFunc:       pflagSetFunc(fs.String, fs.StringP),
		BoolSliceFunc:    pflagSetFunc(fs.BoolSlice, fs.BoolSliceP),
		IntSliceFunc:     pflagSetFunc(fs.IntSlice, fs.IntSliceP),
		Int32SliceFunc:   pflagSetFunc(fs.Int32Slice, fs.Int32SliceP),
		Int64SliceFunc:   pflagSetFunc(fs.Int64Slice, fs.Int64SliceP),
		UintSliceFunc:    pflagSetFunc(fs.UintSlice, fs.UintSliceP),
		Float32SliceFunc: pflagSetFunc(fs.Float32Slice, fs.Float32SliceP),
		Float64SliceFunc: pflagSetFunc(fs.Float64Slice, fs.Float64SliceP),
		StringSliceFunc:  pflagSetFunc(fs.StringSlice, fs.StringSliceP),
		DurationFunc:     pflagSetFunc(fs.Duration, fs.DurationP),
		TimeFunc:         pflagTimeSetFunc(fs.String, fs.StringP),
		CountFunc:        pflagCountSetFunc(fs.CountP),
		AnyFunc:          pflagSetFunc(fs.String, fs.StringP),
	}
}

func PFlagGetConverter(fs *pflag.FlagSet) *DefaultConverter {
	return &DefaultConverter{
		BoolFunc:         fs.GetBool,
		IntFunc:          fs.GetInt,
		Int8Func:         fs.GetInt8,
		Int16Func:        fs.GetInt16,
		Int32Func:        fs.GetInt32,
		Int64Func:        fs.GetInt64,
		UintFunc:         fs.GetUint,
		Uint8Func:        fs.GetUint8,
		Uint16Func:       fs.GetUint16,
		Uint32Func:       fs.GetUint32,
		Uint64Func:       fs.GetUint64,
		Float32Func:      fs.GetFloat32,
		Float64Func:      fs.GetFloat64,
		StringFunc:       fs.GetString,
		BoolSliceFunc:    func(name string) ([]bool, error) { return fs.GetBoolSlice(name) },
		IntSliceFunc:     func(name string) ([]int, error) { return fs.GetIntSlice(name) },
		Int32SliceFunc:   func(name string) ([]int32, error) { return fs.GetInt32Slice(name) },
		Int64SliceFunc:   func(name string) ([]int64, error) { return fs.GetInt64Slice(name) },
		UintSliceFunc:    func(name string) ([]uint, error) { return fs.GetUintSlice(name) },
		Float32SliceFunc: func(name string) ([]float32, error) { return fs.GetFloat32Slice(name) },
		Float64SliceFunc: func(name string) ([]float64, error) { return fs.GetFloat64Slice(name) },
		StringSliceFunc:  func(name string) ([]string, error) { return fs.GetStringSlice(name) },
		DurationFunc:     func(name string) (time.Duration, error) { return fs.GetDuration(name) },
		TimeFunc: func(name string) (time.Time, error) {
			s, err := fs.GetString(name)
			if err != nil {
				return time.Time{}, err
			}
			if s == "" {
				return time.Time{}, nil
			}
			return ParseTime(s)
		},
		CountFunc: func(name string) (int, error) {
			return fs.GetCount(name)
		},
	}
}
