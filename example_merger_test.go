package structconfig_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/berquerant/structconfig"
)

func ExampleMerger() {
	type T struct {
		I                 int    `name:"i" default:"1"`
		S                 string `name:"s" default:"s"`
		II                []int  `name:"ii" default:"[1]"`
		IgnoreWithoutName int
	}

	callback := func(s structconfig.StructField, v string, fv func() reflect.Value) error {
		if s.Name() != "II" {
			return errors.New("unexpected field name")
		}
		var xs []int
		if err := json.Unmarshal([]byte(v), &xs); err != nil {
			return err
		}
		fv().Set(reflect.ValueOf(xs))
		return nil
	}

	eq := func(a, b any) (bool, error) {
		// expect only []int because int and string are supported by structconfig
		xs, ok := a.([]int)
		if !ok {
			return false, nil
		}
		ys, ok := b.([]int)
		if !ok {
			return false, nil
		}
		if len(xs) != len(ys) {
			return false, nil
		}
		for i, x := range xs {
			if x != ys[i] {
				return false, nil
			}
		}
		return true, nil
	}

	m := structconfig.NewMerger[T](
		structconfig.WithAnyCallback(callback),
		structconfig.WithAnyEqual(eq),
	)
	got, err := m.Merge(
		T{
			I:  100,
			S:  "s", // default
			II: []int{100},
		},
		T{
			I:  1, // default
			S:  "win",
			II: []int{1}, // default
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(got.I, got.S, got.II)
	// Output: 100 win [100]
}
