package internal_test

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestDefaultReceptor(t *testing.T) {
	type T struct {
		B         bool    `default:"true"`
		I         int     `default:"1"`
		U         uint    `default:"10"`
		F         float32 `default:"1.1"`
		S         string  `default:"str"`
		NoDefault int
		Slice     []int `default:"[1,2]"`
	}

	want := T{
		B:     true,
		I:     1,
		U:     10,
		F:     1.1,
		S:     "str",
		Slice: []int{1, 2},
	}

	var got T

	r, err := internal.DefaultReceptor(
		&got,
		func(_ internal.StructField, v string, fv func() reflect.Value) error {
			var xs []int
			if err := json.Unmarshal([]byte(v), &xs); err != nil {
				return err
			}
			fv().Set(reflect.ValueOf(xs))
			return nil
		},
	)
	assert.Nil(t, err)

	typ, err := internal.NewType(got, "")
	assert.Nil(t, err)

	assert.Nil(t, typ.Accept(r))
	assert.Equal(t, want, got)
}
