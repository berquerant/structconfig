package internal

import (
	"reflect"

	"github.com/spf13/pflag"
)

func PFlagSetReceptor(fs *pflag.FlagSet) *PairsReceptor {
	return FlagSetReceptor(PFlagSetTypeReceptor(fs))
}

func PFlagGetReceptor(
	ptr any,
	fs *pflag.FlagSet,
	anyCallback func(StructField, string, func() reflect.Value) error,
) (*PairsReceptor, error) {
	typedReceptor, err := SetTypedReceptor(ptr, anyCallback)
	if err != nil {
		return nil, err
	}
	get := func(s StructField) (string, error) {
		if name, ok := s.Tag().Name(); ok {
			return name, nil
		}
		return "", ErrParseAsDefault
	}
	return PairsSynthReceptor(
		get,
		PFlagGetConverter(fs),
		typedReceptor,
	), nil
}

func pflagSetFunc[T any](f func(string, T, string) *T) TypedReceptorFunc[T] {
	return func(s StructField, defaultValue T) error {
		if name, ok := s.Tag().Name(); ok {
			_ = f(name, defaultValue, s.Tag().Usage())
		}
		return nil
	}
}

func PFlagSetTypeReceptor(fs *pflag.FlagSet) *DefaultTypedReceptor {
	return &DefaultTypedReceptor{
		BoolFunc:    pflagSetFunc(fs.Bool),
		IntFunc:     pflagSetFunc(fs.Int),
		Int8Func:    pflagSetFunc(fs.Int8),
		Int16Func:   pflagSetFunc(fs.Int16),
		Int32Func:   pflagSetFunc(fs.Int32),
		Int64Func:   pflagSetFunc(fs.Int64),
		UintFunc:    pflagSetFunc(fs.Uint),
		Uint8Func:   pflagSetFunc(fs.Uint8),
		Uint16Func:  pflagSetFunc(fs.Uint16),
		Uint32Func:  pflagSetFunc(fs.Uint32),
		Uint64Func:  pflagSetFunc(fs.Uint64),
		Float32Func: pflagSetFunc(fs.Float32),
		Float64Func: pflagSetFunc(fs.Float64),
		StringFunc:  pflagSetFunc(fs.String),
		AnyFunc:     pflagSetFunc(fs.String),
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
