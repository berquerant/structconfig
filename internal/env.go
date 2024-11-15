package internal

import (
	"os"
	"regexp"
	"strings"
)

var (
	regexpEnvNameReplace = regexp.MustCompile(`[-.]`)
)

// EnvVar is a environment variable.
type EnvVar string

func NewEnvVar(name string) EnvVar {
	s := regexpEnvNameReplace.ReplaceAllLiteralString(name, "_")
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
