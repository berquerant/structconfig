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
		N  int    `name:"int_value" default:"10"`
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
	// Output: true str 10 0 0
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
		Ignore3 struct {
			S string
		}
	}

	anyCallback := func(s structconfig.StructField, v string, fv func() reflect.Value) error {
		name, ok := s.Tag().Name()
		if !ok {
			// ignore fields without name
			return nil
		}
		switch name {
		case "struct_value":
			fv().Set(reflect.ValueOf(struct {
				S string
			}{
				S: v,
			}))
			return nil
		default:
			return errors.New("unexpected tag name")
		}
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

func Example_stringSliceSplitAndSep() {
	// split: false (default): does not split elements by delimiter.
	// - CLI flags: multiple flags are accumulated as elements (e.g. --raw "a,b" --raw "c" => []string{"a,b", "c"}).
	// - Env / Default: a single string is assigned as a 1-element slice (e.g. RAW="a,b" => []string{"a,b"}).
	//
	// split: true (default sep: ","): splits elements by comma (CSV format).
	// - e.g. --csv "a,b" --csv "c" => []string{"a", "b", "c"}.
	//
	// split: true + sep: ";": splits elements by custom delimiter.
	// - e.g. --custom "a;b" --custom "c" => []string{"a", "b", "c"}.
	type T struct {
		Raw    []string `name:"raw" default:"a,b"`
		Csv    []string `name:"csv" default:"a,b" split:"true"`
		Custom []string `name:"custom" default:"a;b" split:"true" sep:";"`
	}

	sc := structconfig.New[T]()

	// 1. FromDefault
	var def T
	if err := sc.FromDefault(&def); err != nil {
		panic(err)
	}
	fmt.Printf("Default: Raw=%v, Csv=%v, Custom=%v\n", def.Raw, def.Csv, def.Custom)

	// 2. FromEnv
	envs := map[string]string{
		"RAW":    "x,y",
		"CSV":    "x,y",
		"CUSTOM": "x;y",
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	var fromEnv T
	if err := sc.FromEnv(&fromEnv); err != nil {
		panic(err)
	}
	fmt.Printf("Env: Raw=%v, Csv=%v, Custom=%v\n", fromEnv.Raw, fromEnv.Csv, fromEnv.Custom)

	// 3. FromFlags
	fs := pflag.NewFlagSet("example", pflag.ContinueOnError)
	if err := sc.SetFlags(fs); err != nil {
		panic(err)
	}
	args := []string{
		"--raw", "1,2", "--raw", "3",
		"--csv", "1,2", "--csv", "3",
		"--custom", "1;2", "--custom", "3",
	}
	if err := fs.Parse(args); err != nil {
		panic(err)
	}

	var fromFlags T
	if err := sc.FromFlags(&fromFlags, fs); err != nil {
		panic(err)
	}
	fmt.Printf("Flags: Raw=%v, Csv=%v, Custom=%v\n", fromFlags.Raw, fromFlags.Csv, fromFlags.Custom)

	// Output:
	// Default: Raw=[a,b], Csv=[a b], Custom=[a b]
	// Env: Raw=[x,y], Csv=[x y], Custom=[x y]
	// Flags: Raw=[1,2 3], Csv=[1 2 3], Custom=[1 2 3]
}

func ExampleWithEnvPrefix() {
	// WithEnvPrefix adds a prefix only to environment variable names,
	// keeping command-line flag names unchanged.
	type Config struct {
		Host string `name:"host" default:"localhost"`
		Port int    `name:"port" default:"8080"`
	}

	envs := map[string]string{
		"MYAPP_HOST": "127.0.0.1",
		"MYAPP_PORT": "3000",
		"HOST":       "otherhost", // ignored because of prefix
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	// 1. Specify WithEnvPrefix when creating StructConfig
	sc := structconfig.New[Config](structconfig.WithEnvPrefix("MYAPP_"))
	var got Config
	if err := sc.FromEnv(&got); err != nil {
		panic(err)
	}
	fmt.Printf("Env: Host=%s, Port=%d\n", got.Host, got.Port)

	// 2. Command-line flags still use original name tag without prefix ("--host", "--port")
	fs := pflag.NewFlagSet("example", pflag.ContinueOnError)
	if err := sc.SetFlags(fs); err != nil {
		panic(err)
	}
	if err := fs.Parse([]string{"--port", "9000"}); err != nil {
		panic(err)
	}
	var fromFlags Config
	if err := sc.FromFlags(&fromFlags, fs); err != nil {
		panic(err)
	}
	fmt.Printf("Flags: Host=%s, Port=%d\n", fromFlags.Host, fromFlags.Port)

	// Output:
	// Env: Host=127.0.0.1, Port=3000
	// Flags: Host=localhost, Port=9000
}

func Example_prefixVsEnvPrefix() {
	// WithPrefix: specifies a prefix for struct tag keys (e.g. tag "app_name" instead of "name").
	// WithEnvPrefix: specifies a prefix for environment variable names (e.g. "MYAPP_PORT" instead of "PORT").
	//
	// When used together:
	// - StructConfig looks up the tag with prefix: `app_name:"port"` -> tag name is "port"
	// - Flag name uses the tag name: "--port"
	// - Environment variable name prepends WithEnvPrefix: "MYAPP_PORT"
	type Config struct {
		Port int    `app_name:"port" app_default:"8080"`
		Host string `app_name:"host" app_default:"localhost"`
	}

	envs := map[string]string{
		"MYAPP_PORT": "9000",
		"MYAPP_HOST": "127.0.0.1",
		"PORT":       "3000", // ignored (missing MYAPP_ prefix)
	}
	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	sc := structconfig.New[Config](
		structconfig.WithPrefix("app_"),      // reads `app_name` and `app_default` tags
		structconfig.WithEnvPrefix("MYAPP_"), // reads `MYAPP_<name>` env vars
	)

	// 1. FromEnv reads MYAPP_PORT and MYAPP_HOST
	var fromEnv Config
	if err := sc.FromEnv(&fromEnv); err != nil {
		panic(err)
	}
	fmt.Printf("Env: Host=%s, Port=%d\n", fromEnv.Host, fromEnv.Port)

	// 2. Flags are registered with the tag name "--port", "--host" (unaffected by WithEnvPrefix)
	fs := pflag.NewFlagSet("example", pflag.ContinueOnError)
	if err := sc.SetFlags(fs); err != nil {
		panic(err)
	}
	if err := fs.Parse([]string{"--port", "7070"}); err != nil {
		panic(err)
	}
	var fromFlags Config
	if err := sc.FromFlags(&fromFlags, fs); err != nil {
		panic(err)
	}
	fmt.Printf("Flags: Host=%s, Port=%d\n", fromFlags.Host, fromFlags.Port)

	// Output:
	// Env: Host=127.0.0.1, Port=9000
	// Flags: Host=localhost, Port=7070
}
