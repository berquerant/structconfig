package internal

type FlagSetFunc[T any] func(name string, defaultValue T, usage string) error

func (f FlagSetFunc[T]) SetFlag(s StructField, v T) error {
	if f == nil {
		return nil
	}

	name, ok := s.Tag().Name()
	if !ok {
		// ignore the field
		return nil
	}
	usage := s.Tag().Usage()
	return f(name, v, usage)
}

func FlagSetReceptor(typedReceptor TypedReceptor) *PairsReceptor {
	get := func(s StructField) (string, error) {
		if v, ok := s.Tag().Default(); ok {
			return v, nil
		}
		return "", ErrParseAsDefault
	}
	return PairsSynthReceptor(get, NewConv(), typedReceptor)
}

func FlagSetTypedReceptor(
	boolFunc FlagSetFunc[bool],
	intFunc FlagSetFunc[int],
	int8Func FlagSetFunc[int8],
	int16Func FlagSetFunc[int16],
	int32Func FlagSetFunc[int32],
	int64Func FlagSetFunc[int64],
	uintFunc FlagSetFunc[uint],
	uint8Func FlagSetFunc[uint8],
	uint16Func FlagSetFunc[uint16],
	uint32Func FlagSetFunc[uint32],
	uint64Func FlagSetFunc[uint64],
	float32Func FlagSetFunc[float32],
	float64Func FlagSetFunc[float64],
	stringFunc FlagSetFunc[string],
	anyFunc FlagSetFunc[string],
) *DefaultTypedReceptor {
	return &DefaultTypedReceptor{
		BoolFunc:    boolFunc.SetFlag,
		IntFunc:     intFunc.SetFlag,
		Int8Func:    int8Func.SetFlag,
		Int16Func:   int16Func.SetFlag,
		Int32Func:   int32Func.SetFlag,
		Int64Func:   int64Func.SetFlag,
		UintFunc:    uintFunc.SetFlag,
		Uint8Func:   uint8Func.SetFlag,
		Uint16Func:  uint16Func.SetFlag,
		Uint32Func:  uint32Func.SetFlag,
		Uint64Func:  uint64Func.SetFlag,
		Float32Func: float32Func.SetFlag,
		Float64Func: float64Func.SetFlag,
		StringFunc:  stringFunc.SetFlag,
		AnyFunc:     anyFunc.SetFlag,
	}
}
