package structconfig_test

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

func ExampleBuilder() {
	type T struct {
		Default  string `json:"default" name:"default_value" default:"default"`
		Env      string `json:"env" name:"env_value" default:"env_default"`
		Flag     string `json:"flag" name:"flag_value" default:"flag_default"`
		File     string `json:"file" name:"file_value" default:"file_default"`
		Override string `json:"override" name:"override_value" default:"override_default"`
	}
	file := []byte(`{"file":"from_file","override":"from_file"}`)
	envs := map[string]string{
		"ENV_VALUE":      "from_env",
		"OVERRIDE_VALUE": "should_not_appear",
	}
	args := []string{
		"--flag_value", "from_flag",
		"--override_value", "overrided",
	}

	for k, v := range envs {
		os.Setenv(k, v)
	}
	defer func() {
		for k := range envs {
			os.Unsetenv(k)
		}
	}()

	fs := pflag.NewFlagSet("test", pflag.ContinueOnError)
	c, err := structconfig.NewBuilder[T](structconfig.New[T](), structconfig.NewMerger[T]()).
		Add(func(sc *structconfig.StructConfig[T]) (*T, error) {
			var t T
			if err := sc.FromDefault(&t); err != nil {
				return nil, err
			}
			if err := json.Unmarshal(file, &t); err != nil {
				return nil, err
			}
			return &t, nil
		}).
		Add(func(sc *structconfig.StructConfig[T]) (*T, error) {
			var t T
			if err := sc.FromEnv(&t); err != nil {
				return nil, err
			}
			return &t, nil
		}).
		Add(func(sc *structconfig.StructConfig[T]) (*T, error) {
			if err := sc.SetFlags(fs); err != nil {
				return nil, err
			}
			if err := fs.Parse(args); err != nil {
				return nil, err
			}
			var t T
			if err := sc.FromFlags(&t, fs); err != nil {
				return nil, err
			}
			return &t, nil
		}).
		Build()
	if err != nil {
		panic(err)
	}
	fmt.Println(c.Default, c.Env, c.Flag, c.File, c.Override)
	// Output: default from_env from_flag from_file overrided
}
