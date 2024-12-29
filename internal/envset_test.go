package internal_test

import (
	"encoding/json"
	"os"
	"reflect"
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestEnvReceptor(t *testing.T) {
	type T struct {
		B         bool    `name:"eb"`
		I         int     `name:"ei"`
		U         uint    `name:"eu"`
		F         float32 `name:"ef"`
		S         string  `name:"es" default:"str"`
		NoDefault int     `name:"eno_default"`
		Slice     []int   `name:"eslice"`
		Ignore    []int   `name:"-"`
	}

	envs := map[string]string{
		"EB": "true",
		"EI": "1",
		"EU": "10",
		"EF": "1.1",
		// ENO_DEFAULT is not defined
		"ESLICE": "[1,2]",
		"IGNORE": "[1]",
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	want := T{
		B:     true,
		I:     1,
		U:     10,
		F:     1.1,
		S:     "str",
		Slice: []int{1, 2},
	}

	var got T

	r, err := internal.EnvReceptor(
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
