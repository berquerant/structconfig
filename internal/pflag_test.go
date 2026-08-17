package internal_test

import (
	"reflect"
	"sort"
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

func TestPFlagReceptor(t *testing.T) {
	type CustomStruct struct {
		Name string
	}

	type T struct {
		B         bool         `name:"fb" default:"false" usage:"BOOL"`
		I         int          `name:"fi" default:"1" usage:"INT"`
		U         uint         `name:"fu" default:"10" usage:"UINT"`
		F         float32      `name:"ff" default:"1.1" usage:"FLOAT"`
		S         string       `name:"fs" default:"str" usage:"STRING" short:"s"`
		NoDefault int          `name:"fnodefault"`
		V         CustomStruct `name:"fv" default:"def"`
		List1     []string     `name:"flist1"`
		List2     []string     `name:"flist2" split:"true"`
		List3     []string     `name:"flist3" split:"true" sep:";"`
		Ignore1   int
		Ignore2   int `name:"-" default:"1000" usage:"IGNORE2"`
	}

	flagNames := []string{
		"fb",
		"fi",
		"fu",
		"ff",
		"fs",
		"fnodefault",
		"fv",
		"flist1",
		"flist2",
		"flist3",
	}
	sort.Strings(flagNames)

	var typVal T
	typ, err := internal.NewType(typVal, "")
	assert.Nil(t, err)

	for _, tc := range []struct {
		title string
		args  []string
		want  T
	}{
		{
			title: "shorthand",
			args: []string{
				"-s", "SHORT",
			},
			want: T{
				I: 1,
				U: 10,
				F: 1.1,
				S: "SHORT",
				V: CustomStruct{Name: "def"},
			},
		},
		{
			title: "slice split and nosplit",
			args: []string{
				"--flist1", "a,b",
				"--flist1", "c",
				"--flist2", "a,b",
				"--flist2", "c",
				"--flist3", "a;b",
				"--flist3", "c",
			},
			want: T{
				I:     1,
				U:     10,
				F:     1.1,
				S:     "str",
				V:     CustomStruct{Name: "def"},
				List1: []string{"a,b", "c"},
				List2: []string{"a", "b", "c"},
				List3: []string{"a", "b", "c"},
			},
		},
		{
			title: "change all",
			args: []string{
				"--fb",
				"--fi", "-2",
				"--fu", "3",
				"--ff", "10.1",
				"--fs", "changed",
				"--fnodefault", "-24",
				"--fv", "custom",
			},
			want: T{
				B:         true,
				I:         -2,
				U:         3,
				F:         10.1,
				S:         "changed",
				NoDefault: -24,
				V:         CustomStruct{Name: "custom"},
			},
		},
		{
			title: "change fv",
			args:  []string{"--fv", "custom_only"},
			want: T{
				I: 1,
				U: 10,
				F: 1.1,
				S: "str",
				V: CustomStruct{Name: "custom_only"},
			},
		},
		{
			title: "change fi",
			args:  []string{"--fi", "1000"},
			want: T{
				I: 1000,
				U: 10,
				F: 1.1,
				S: "str",
				V: CustomStruct{Name: "def"},
			},
		},
		{
			title: "default",
			want: T{
				I: 1,
				U: 10,
				F: 1.1,
				S: "str",
				V: CustomStruct{Name: "def"},
			},
		},
	} {
		t.Run(tc.title, func(t *testing.T) {
			fs := pflag.NewFlagSet("test", pflag.ContinueOnError)

			t.Run("Set", func(t *testing.T) {
				r := internal.PFlagSetReceptor(fs)
				assert.Nil(t, typ.Accept(r))
			})

			t.Run("CheckFlagNames", func(t *testing.T) {
				names := []string{}
				fs.VisitAll(func(f *pflag.Flag) {
					names = append(names, f.Name)
				})
				sort.Strings(names)
				assert.Equal(t, flagNames, names)
			})

			t.Run("Parse", func(t *testing.T) {
				assert.Nil(t, fs.Parse(tc.args))
			})

			t.Run("Get", func(t *testing.T) {
				var got T
				r, err := internal.PFlagGetReceptor(
					&got,
					fs,
					func(_ internal.StructField, v string, fv func() reflect.Value) error {
						fv().Set(reflect.ValueOf(CustomStruct{Name: v}))
						return nil
					},
					nil,
				)
				assert.Nil(t, err)

				assert.Nil(t, typ.Accept(r))
				assert.Equal(t, tc.want, got)
			})
		})
	}

}
