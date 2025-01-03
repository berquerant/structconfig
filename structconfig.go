package structconfig

import (
	"reflect"

	"github.com/berquerant/structconfig/internal"
	"github.com/spf13/pflag"
)

const (
	TagName    = internal.TagName
	TagUsage   = internal.TagUsage
	TagDefault = internal.TagDefault
	TagShort   = internal.TagShort
)

var (
	ErrStructConfig     = internal.ErrStructConfig
	ErrNotStruct        = internal.ErrNotStruct
	ErrNotStructPointer = internal.ErrNotStructPointer
)

type (
	StructField     = internal.StructField
	Type            = internal.Type
	Tag             = internal.Tag
	EnvVar          = internal.EnvVar
	Receptor        = internal.Receptor
	BoolReceptor    = internal.BoolReceptor
	IntReceptor     = internal.IntReceptor
	UintReceptor    = internal.UintReceptor
	FloatReceptor   = internal.FloatReceptor
	StringReceptor  = internal.StringReceptor
	AnyReceptor     = internal.AnyReceptor
	AnyCallbackFunc = func(StructField, string, func() reflect.Value) error
	AnyEqualFunc    = func(left, right any) (bool, error)
	Unsigned        = internal.Unsigned
	Supported       = internal.Supported
)

func IsSupportedKind(k reflect.Kind) bool         { return internal.IsSupportedKind(k) }
func NewType(v any, prefix string) (*Type, error) { return internal.NewType(v, prefix) }

//go:generate go run github.com/berquerant/goconfig -configOption Option -option -output structconfig_config_generated.go -field "AnyCallback AnyCallbackFunc|AnyEqual AnyEqualFunc|Prefix string"

func newDefaultConfigBuilder() *ConfigBuilder {
	return NewConfigBuilder().
		AnyCallback(nil).
		AnyEqual(nil).
		Prefix("")
}

type Merger[T any] struct {
	*internal.Merger[T]
}

// NewMerger returns a new Merger.
//
// AnyCallback parses "default" tag value and set it.
// AnyEqual reports true if left equals right when kind of arguments are not supported.
// Prefix adds a prefix to "default" tag name.
func NewMerger[T any](opt ...Option) *Merger[T] {
	c := newDefaultConfigBuilder().Build()
	c.Apply(opt...)

	return &Merger[T]{
		internal.NewMerger[T](
			c.AnyCallback.Get(),
			c.AnyEqual.Get(),
			c.Prefix.Get(),
		),
	}
}

// Merge values based on the 'default' tag values.
// For each field with 'name' and 'default' tags, if the right value is not the default, use it; if not, use the left value.
// If that is also the default, set the default value. Return this instance.
func (m *Merger[T]) Merge(left, right T) (T, error) {
	return m.Merger.Merge(left, right)
}

// New returns a new StructConfig.
//
// AnyCallback parses "default" tag value and set it.
// Prefix adds a prefix to "name", "short", "default" and "usage" tag name.
func New[T any](opt ...Option) *StructConfig[T] {
	c := newDefaultConfigBuilder().Build()
	c.Apply(opt...)

	return &StructConfig[T]{
		anyCallback: c.AnyCallback.Get(),
		prefix:      c.Prefix.Get(),
	}
}

type StructConfig[T any] struct {
	anyCallback AnyCallbackFunc
	prefix      string
}

func (sc StructConfig[T]) newType() (*Type, error) {
	var t T
	return NewType(t, sc.prefix)
}

func (sc StructConfig[T]) from(r Receptor) error {
	typ, err := sc.newType()
	if err != nil {
		return err
	}
	return typ.Accept(r)
}

// FromDefault sets "default" tag values to v.
func (sc StructConfig[T]) FromDefault(v *T) error {
	r, err := internal.DefaultReceptor(v, sc.anyCallback)
	if err != nil {
		return err
	}
	return sc.from(r)
}

// FromEnv sets environment variable values to v.
//
// Environment variable name will be
//
//	NewEnvVar("name tag value").String()
//
// All '.' and '-' will be replaced with '_', making it all uppsercase.
func (sc StructConfig[T]) FromEnv(v *T) error {
	r, err := internal.EnvReceptor(v, sc.anyCallback)
	if err != nil {
		return err
	}
	return sc.from(r)
}

// FromFlags sets values to v from command-line flags.
//
// Flag name is from "name" tag value.
func (sc StructConfig[T]) FromFlags(v *T, fs *pflag.FlagSet) error {
	r, err := internal.PFlagGetReceptor(v, fs, sc.anyCallback)
	if err != nil {
		return err
	}
	return sc.from(r)
}

// SetFlags sets command-line flags.
//
// Flag name is from "name" tag value.
// Flag shorthand is from "short" tag value.
// Flag default value is from "default" tag value.
// Flag usage is from "usage" tag value.
func (sc StructConfig[T]) SetFlags(fs *pflag.FlagSet) error {
	r := internal.PFlagSetReceptor(fs)
	return sc.from(r)
}
