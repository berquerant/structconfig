[![Go Reference](https://pkg.go.dev/badge/github.com/berquerant/structconfig.svg)](https://pkg.go.dev/github.com/berquerant/structconfig)

# structconfig

Map default values, environment variables, and command-line arguments to struct tags.

# Installation

``` shell
go get github.com/berquerant/structconfig
```

# Examples

## Default values

``` go
type T struct {
  I int `default:"10"`
}

sc := structconfig.New[T]()
var got T
if err := sc.FromDefault(&got); err != nil {
  panic(err)
}
// got.I == 10
```

## Environment variables

``` go
type T struct {
  I int `name:"int_value"`
}

os.Setenv("INT_VALUE", "10")
sc := structconfig.New[T]()
var got T
if err := sc.FromEnv(&got); err != nil {
  panic(err)
}
// got.I == 10
```

## Command-line flags ([pflag](https://github.com/spf13/pflag))

``` go
type T struct {
  I int `name:"int_value" default:"10"`
}

var fs *pflag.FlagSet = // ...
sc := structconfig.New[T]()
if err := sc.SetFlags(fs); err != nil {
  panic(err)
}
if err := fs.Parse([]string{"--int_value", "100"}); err != nil {
  panic(err)
}
var got T
if err := sc.FromFlags(&got, fs); err != nil {
  panic(err)
}
// got.I == 100
```

## Structured Logging

You can pass a `*slog.Logger` to log configuration loading and merge decisions at **Debug** level:

``` go
logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
sc := structconfig.New[T](structconfig.WithLogger(logger))
// Logs at Debug level when values are set from default, env, or flag

merger := structconfig.NewMerger[T](structconfig.WithLogger(logger))
// Logs at Debug level when fields are merged/overridden
```

## More examples

- [Merger](example_merger_test.go)
- [Default, Env, Flag](example_structconfig_test.go)
