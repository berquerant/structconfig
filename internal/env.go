package internal

import (
	"os"
	"regexp"
	"strings"
)

var (
	regexpEnvNameReplace = regexp.MustCompile(`[-.]`)
)

type EnvVar string

func NewEnvVar(name string) EnvVar {
	s := regexpEnvNameReplace.ReplaceAllLiteralString(name, "_")
	return EnvVar(strings.ToUpper(s))
}

func (v EnvVar) Get() (string, bool) {
	return os.LookupEnv(string(v))
}

func (v EnvVar) String() string {
	return string(v)
}
