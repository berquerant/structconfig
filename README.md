# structconfig

Map default values, environment variables, and command-line arguments to struct tags.

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

## More examples

- [Merger](example_merger_test.go)
- [Default, Env, Flag](example_structconfig_test.go)
