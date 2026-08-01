package internal

import "reflect"

// SetReceptor returns a new PairsReceptor to set value to ptr.
//
// ptr should be a pointer of struct.
// get should return a value to be set; false means not found.
// converter should convert result of get().
//
// You can customize how unsupported field values are set by anyCallback,
// if nil, disable this feature.
// Example of anyCallback:
//
//	func(field StructField, str string, fieldPtr func() reflect.Value) error {
//	  var xs []int
//	  if err := json.Unmarshal([]byte(str), &xs); err != nil {
//	    return err
//	  }
//	  fieldPtr().Set(reflect.ValueOf(xs))
//	  return nil
//	}
//
// anyCallback should parse str and set the value to fieldPtr.
func SetReceptor(
	ptr any,
	get func(StructField) (string, error),
	converter Converter,
	anyCallback func(StructField, string, func() reflect.Value) error,
	onSet func(StructField, any),
) (*PairsReceptor, error) {
	typedReceptor, err := SetTypedReceptor(ptr, anyCallback, onSet)
	if err != nil {
		return nil, err
	}
	return PairsSynthReceptor(get, converter, typedReceptor), nil
}

func SetTypedReceptor(
	ptr any,
	anyCallback func(StructField, string, func() reflect.Value) error,
	onSet func(StructField, any),
) (*DefaultTypedReceptor, error) {
	typ := reflect.TypeOf(ptr)
	if !(typ.Kind() == reflect.Pointer && typ.Elem().Kind() == reflect.Struct) {
		return nil, JoinErrors(ErrNotStructPointer)
	}

	v := reflect.ValueOf(ptr)
	fv := func(s StructField) reflect.Value {
		return v.Elem().FieldByName(s.Name())
	}
	notify := func(s StructField, val any) {
		if onSet != nil {
			onSet(s, val)
		}
	}

	return &DefaultTypedReceptor{
		BoolFunc: func(s StructField, v bool) error {
			fv(s).SetBool(v)
			notify(s, v)
			return nil
		},
		IntFunc: func(s StructField, v int) error {
			fv(s).SetInt(int64(v))
			notify(s, v)
			return nil
		},
		Int8Func: func(s StructField, v int8) error {
			fv(s).SetInt(int64(v))
			notify(s, v)
			return nil
		},
		Int16Func: func(s StructField, v int16) error {
			fv(s).SetInt(int64(v))
			notify(s, v)
			return nil
		},
		Int32Func: func(s StructField, v int32) error {
			fv(s).SetInt(int64(v))
			notify(s, v)
			return nil
		},
		Int64Func: func(s StructField, v int64) error {
			fv(s).SetInt(v)
			notify(s, v)
			return nil
		},
		UintFunc: func(s StructField, v uint) error {
			fv(s).SetUint(uint64(v))
			notify(s, v)
			return nil
		},
		Uint8Func: func(s StructField, v uint8) error {
			fv(s).SetUint(uint64(v))
			notify(s, v)
			return nil
		},
		Uint16Func: func(s StructField, v uint16) error {
			fv(s).SetUint(uint64(v))
			notify(s, v)
			return nil
		},
		Uint32Func: func(s StructField, v uint32) error {
			fv(s).SetUint(uint64(v))
			notify(s, v)
			return nil
		},
		Uint64Func: func(s StructField, v uint64) error {
			fv(s).SetUint(v)
			notify(s, v)
			return nil
		},
		Float32Func: func(s StructField, v float32) error {
			fv(s).SetFloat(float64(v))
			notify(s, v)
			return nil
		},
		Float64Func: func(s StructField, v float64) error {
			fv(s).SetFloat(v)
			notify(s, v)
			return nil
		},
		StringFunc: func(s StructField, v string) error {
			fv(s).SetString(v)
			notify(s, v)
			return nil
		},
		AnyFunc: func(s StructField, v string) error {
			if anyCallback == nil {
				return nil
			}
			err := anyCallback(
				s,
				v,
				func() reflect.Value { return fv(s) },
			)
			if err == nil {
				notify(s, fv(s).Interface())
			}
			return err
		},
	}, nil
}
