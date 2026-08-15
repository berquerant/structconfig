package internal

import (
	"os"
	"strings"
)

var (
	envNameReplacer = strings.NewReplacer("-", "_", ".", "_")
)

// EnvVar is a environment variable.
type EnvVar string

func NewEnvVar(name string) EnvVar {
	s := envNameReplacer.Replace(name)
	return EnvVar(strings.ToUpper(s))
}

// Get retrieves the value of the environment variable.
func (v EnvVar) Get() (string, bool) {
	return os.LookupEnv(string(v))
}

// String returns the name of the environment variable.
func (v EnvVar) String() string {
	return string(v)
}
