package internal_test

import (
	"testing"

	"github.com/berquerant/structconfig/internal"
	"github.com/stretchr/testify/assert"
)

func TestEnvVar(t *testing.T) {
	for _, tc := range []struct {
		name string
		want string
	}{
		{
			name: "a",
			want: "A",
		},
		{
			name: "A",
			want: "A",
		},
		{
			name: "ENV_VAR",
			want: "ENV_VAR",
		},
		{
			name: "env-var",
			want: "ENV_VAR",
		},
		{
			name: "Env.Var",
			want: "ENV_VAR",
		},
		{
			name: "Env.Var.x1-Z",
			want: "ENV_VAR_X1_Z",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.want, internal.NewEnvVar(tc.name).String())
		})
	}
}
