package structconfig_test

import (
	"fmt"
	"os"

	"github.com/berquerant/structconfig"
	"github.com/spf13/pflag"
)

func ExampleNewConfigWithMerge() {
	type T struct {
		Default  string `name:"default_value" default:"default"`
		Env      string `name:"env_value" default:"env_default"`
		Flag     string `name:"flag_value" default:"flag_default"`
		Override string `name:"override_value" default:"override_default"`
	}
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
	c, err := structconfig.NewConfigWithMerge[T](
		structconfig.New[T](),
		structconfig.NewMerger[T](),
		fs,
		structconfig.WithArguments(args),
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(c.Default, c.Env, c.Flag, c.Override)
	// Output: default from_env from_flag overrided
}
