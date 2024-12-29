package structconfig_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

func ExampleStructConfig_FromDefault() {
	type T struct {
		B bool `default:"false"`
		I int
		F float32 `default:"1.1"`
		S string  `default:"str"`
	}

	sc := structconfig.New[T]()
	var got T
	if err := sc.FromDefault(&got); err != nil {
		panic(err)
	}
	fmt.Println(got.B, got.I, got.F, got.S)
	// Output: false 0 1.1 str
}

func ExampleStructConfig_FromEnv() {
	type T struct {
		B  bool   `name:"bool_value"`
		S  string `name:"string_value"`
		N  int    `name:"int_value"`
		N2 int
		N3 int `name:"-"`
	}

	envs := map[string]string{
		"BOOL_VALUE":   "true",
		"STRING_VALUE": "str",
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	sc := structconfig.New[T]()
	var got T
	if err := sc.FromEnv(&got); err != nil {
		panic(err)
	}
	fmt.Println(got.B, got.S, got.N, got.N2, got.N3)
	// Output: true str 0 0 0
}

func ExampleStructConfig_FromFlags() {
	type T struct {
		B       bool   `name:"bool_value" usage:"BOOL"`
		S       string `name:"string_value" default:"str"`
		X       bool   `name:"bool_short" short:"x"`
		Ignore1 int
		Ignore2 int `name:"-"`
		V       struct {
			S string
		} `name:"struct_value"`
	}

	anyCallback := func(s structconfig.StructField, v string, fv func() reflect.Value) error {
		if x, ok := s.Tag().Name(); !ok || x != "struct_value" {
			// among T, only struct_value is not supported
			return errors.New("unexpected tag name")
		}
		fv().Set(reflect.ValueOf(struct {
			S string
		}{
			S: v,
		}))
		return nil
	}

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	sc := structconfig.New[T](structconfig.WithAnyCallback(anyCallback))

	if err := sc.SetFlags(fs); err != nil {
		panic(err)
	}

	flagNames := []string{}
	fs.VisitAll(func(f *pflag.Flag) {
		flagNames = append(flagNames, f.Name)
	})

	if err := fs.Parse([]string{"--bool_value", "--struct_value", "sv", "-x"}); err != nil {
		panic(err)
	}

	var got T
	if err := sc.FromFlags(&got, fs); err != nil {
		panic(err)
	}

	sort.Strings(flagNames)
	fmt.Println(strings.Join(flagNames, ","))
	fmt.Println(got.B, got.S, got.Ignore1, got.Ignore2, got.V.S, got.X)
	// Output:
	// bool_short,bool_value,string_value,struct_value
	// true str 0 0 sv true
}

func ExampleStructConfig_FromFlags_prefix() {
	tagPrefix := "sc"
	type T struct {
		B bool   `scname:"bool_value" scusage:"BOOL"`
		S string `scname:"string_value" scdefault:"str"`
	}

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	sc := structconfig.New[T](structconfig.WithPrefix(tagPrefix))

	if err := sc.SetFlags(fs); err != nil {
		panic(err)
	}

	flagNames := []string{}
	fs.VisitAll(func(f *pflag.Flag) {
		flagNames = append(flagNames, f.Name)
	})

	if err := fs.Parse([]string{"--bool_value", "--string_value", "sv"}); err != nil {
		panic(err)
	}

	var got T
	if err := sc.FromFlags(&got, fs); err != nil {
		panic(err)
	}

	sort.Strings(flagNames)
	fmt.Println(strings.Join(flagNames, ","))
	fmt.Println(got.B, got.S)
	// Output:
	// bool_value,string_value
	// true sv
}
