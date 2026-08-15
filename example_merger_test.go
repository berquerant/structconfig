package structconfig_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/berquerant/structconfig"
)

func ExampleMerger() {
	type V struct {
		N int
	}

	type T struct {
		I                 int    `name:"i" default:"1"`
		S                 string `name:"s" default:"s"`
		V                 V      `name:"v" default:"{\"N\":1}"`
		IgnoreWithoutName int
	}

	callback := func(s structconfig.StructField, v string, fv func() reflect.Value) error {
		if s.Name() != "V" {
			return errors.New("unexpected field name")
		}
		var val V
		if err := json.Unmarshal([]byte(v), &val); err != nil {
			return err
		}
		fv().Set(reflect.ValueOf(val))
		return nil
	}

	eq := func(a, b any) (bool, error) {
		// expect only V because int and string are supported by structconfig
		va, ok := a.(V)
		if !ok {
			return false, nil
		}
		vb, ok := b.(V)
		if !ok {
			return false, nil
		}
		return va == vb, nil
	}

	m := structconfig.NewMerger[T](
		structconfig.WithAnyCallback(callback),
		structconfig.WithAnyEqual(eq),
	)
	got, err := m.Merge(
		T{
			I: 100,
			S: "s", // default
			V: V{N: 100},
		},
		T{
			I: 1, // default
			S: "win",
			V: V{N: 1}, // default
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(got.I, got.S, got.V.N)
	// Output: 100 win 100
}
