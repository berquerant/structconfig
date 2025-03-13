package structconfig

import (
	"os"

	"github.com/spf13/pflag"
)

// NewConfigWithMerge generates a Config by taking into account default values,
// environment variables, and command-line arguments.
//
// It overrides the default values with values obtained from environment variables
// and further overrides them with the values from command-line arguments.
// The command-line arguments are obtained by calling [StructConfig.SetFlags]
// on the fs and then parsing with [pflag.FlagSet.Parse].
//
// In this process, the values to be parsed are those specified with [WithArguments].
// If none are specified, it uses [os.Args].
func NewConfigWithMerge[T any](
	sc *StructConfig[T],
	merger *Merger[T],
	fs *pflag.FlagSet,
	opt ...Option,
) (*T, error) {
	c := newDefaultConfigBuilder().Build()
	c.Apply(opt...)

	return NewBuilder(sc, merger).
		Add(func(sc *StructConfig[T]) (*T, error) {
			var t T
			if err := sc.FromEnv(&t); err != nil {
				return nil, err
			}
			return &t, nil
		}).
		Add(func(sc *StructConfig[T]) (*T, error) {
			if err := sc.SetFlags(fs); err != nil {
				return nil, err
			}
			var arguments []string
			if c.Arguments.IsModified() {
				arguments = c.Arguments.Get()
			} else {
				arguments = os.Args
			}
			if err := fs.Parse(arguments); err != nil {
				return nil, err
			}
			var t T
			if err := sc.FromFlags(&t, fs); err != nil {
				return nil, err
			}
			return &t, nil
		}).
		Build()
}
