package internal_test

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestMerger(t *testing.T) {
	var errSome = errors.New("Some")

	type T struct {
		I  int    `default:"1"`
		S  string `default:"s"`
		II []int  `default:"[1]"`
	}
	defaultValue := T{
		I:  1,
		S:  "s",
		II: []int{1},
	}

	callback := func(s internal.StructField, v string, fv func() reflect.Value) error {
		if s.Name() != "II" {
			return errSome
		}
		var xs []int
		if err := json.Unmarshal([]byte(v), &xs); err != nil {
			return err
		}
		fv().Set(reflect.ValueOf(xs))
		return nil
	}
	eq := func(a, b any) (bool, error) {
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

	m := internal.NewMerger[T](
		callback,
		eq,
		"",
	)

	for _, tc := range []struct {
		title string
		left  T
		right T
		want  T
	}{
		{
			title: "merge",
			left: T{
				I:  100,
				S:  "s", // default
				II: []int{100},
			},
			right: T{
				I:  1, // default
				S:  "win",
				II: []int{1}, // default
			},
			want: T{
				I:  100,
				S:  "win",
				II: []int{100},
			},
		},
		{
			title: "all right wins",
			left: T{
				I:  100,
				S:  "lost",
				II: []int{100},
			},
			right: T{
				I:  1000,
				S:  "win",
				II: []int{1000},
			},
			want: T{
				I:  1000,
				S:  "win",
				II: []int{1000},
			},
		},
		{
			title: "right wins",
			left: T{
				S: "lost",
			},
			right: T{
				S: "changed",
			},
			want: T{
				S: "changed",
			},
		},
		{
			title: "right",
			left:  defaultValue,
			right: T{
				S: "changed",
			},
			want: T{
				S: "changed",
			},
		},
		{
			title: "left",
			left: T{
				I: 100,
			},
			right: defaultValue,
			want: T{
				I: 100,
			},
		},
		{
			title: "all default",
			left:  defaultValue,
			right: defaultValue,
			want:  defaultValue,
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			got, err := m.Merge(tc.left, tc.right)
			assert.Nil(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
